package Service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"image/color"
	"log"
	model "main/Model"
	repository "main/Repository"
	"net/http"
	"net/smtp"
	"os"
	"strings"
	"text/template"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"gopkg.in/gomail.v2"

	"github.com/skip2/go-qrcode"
)

var smtpHost = "smtp.gmail.com"
var smtpPort = "587"

type EmailParams struct {
	Subject string
	Body    string
}

func sendEmail(from string, to []string, auth smtp.Auth, params EmailParams) error {
	message := fmt.Sprintf("Subject: %s\n%s", params.Subject, params.Body)
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, []byte(message))
	if err != nil {
		log.Printf("Error sending email: %v\n", err)
		return err
	}
	log.Println("Email Sent!")
	return nil
}

func SetMailSender(email string) (string, []string, smtp.Auth) {
	from := "nemanja.ranitovic@fabricus.tech"
	password, passErr := os.LookupEnv("GOOGLE_MAIL_PASSWORD")
	if !passErr {
		log.Fatal("Google_mail_password not declared in .env file!")
	}
	to := []string{email}
	auth := smtp.PlainAuth("", from, password, smtpHost)
	return from, to, auth
}
func LoadHTMLTemplate(filePath string) (*template.Template, error) {
	tmpl, err := template.ParseFiles(filePath)
	if err != nil {
		return nil, err
	}
	return tmpl, nil
}

func generateLink(baseLink string, replacements map[string]string) string {
	for placeholder, value := range replacements {
		baseLink = strings.Replace(baseLink, placeholder, value, 1)
	}
	return baseLink
}

func SendEmailWithHTMLTemplate(email, subject, filePath string, data interface{}) error {
	from, to, auth := SetMailSender(email)

	tmpl, err := LoadHTMLTemplate(filePath)
	if err != nil {
		fmt.Println("Error loading HTML template:", err)
		return err
	}

	var body bytes.Buffer
	err = tmpl.Execute(&body, data)
	if err != nil {
		fmt.Println("Error executing template:", err)
		return err
	}

	mimeHeaders := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"

	fullBody := mimeHeaders + body.String()

	sendErr := sendEmail(from, to, auth, EmailParams{Subject: subject, Body: fullBody})
	return sendErr
}

func SendConfirmationMail(userID string, email string) bool {
	replacements := map[string]string{
		"{userID}": userID,
	}

	link := generateLink("http://localhost:8080/users/verify/{userID}", replacements)
	sendErr := SendEmailWithHTMLTemplate(email, "Application confirmation", "Controller/pages/Verify.html", struct{ Message string }{Message: link})
	if sendErr != nil {
		return false
	} else {
		return true
	}
}

// QR code logic
func SendReservationMail(email string, qrFilePath string) bool {
	emailBody := `
	<!DOCTYPE html>
	<html>
		<body>
			<p>Dear User,</p>
			<p>Thank you for your reservation. Please find the attached QR code for your reservation.</p>
		</body>
	</html>
	`

	sendErr := SendEmailWithAttachment(email, "Reservation Confirmation", emailBody, qrFilePath)
	if sendErr != nil {
		fmt.Println("Failed to send email:", sendErr)
		return false
	} else {
		return true
	}
}

func SendEmailWithAttachment(to string, subject string, body string, attachmentPath string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", "nemanja.ranitovic@fabricus.tech")
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	m.Attach(attachmentPath)

	d := gomail.NewDialer("smtp.gmail.com", 587, "nemanja.ranitovic@fabricus.tech", "kudv wdvx ikeb kqru")

	if err := d.DialAndSend(m); err != nil {
		return err
	}
	return nil
}

func GenerateQRCode(reservation model.Reservation, filePath string) error {
	data, err := json.Marshal(reservation)
	if err != nil {
		return fmt.Errorf("failed to marshal reservation: %v", err)
	}

	qr, err := qrcode.New(string(data), qrcode.Medium)
	if err != nil {
		return fmt.Errorf("failed to generate QR code: %v", err)
	}

	qr.BackgroundColor = color.White
	qr.ForegroundColor = color.Black

	if err := qr.WriteFile(256, filePath); err != nil {
		return fmt.Errorf("failed to write QR code to file: %v", err)
	}

	return nil
}

func SendReportReplay(reportID primitive.ObjectID, res http.ResponseWriter) error {

	reportsCollection := repository.GetClient().Database("MedicinskaOprema").Collection("Reports")

	var report model.Report
	err := reportsCollection.FindOne(context.TODO(), bson.M{"_id": reportID}).Decode(&report)
	if err != nil {
		return fmt.Errorf("error fetching report: %v", err)
	}

	user, _ := repository.FindUserById(report.UserID, res)

	userEmail := user.Email
	if userEmail == "" {
		return fmt.Errorf("user email not found in report")
	}

	emailData := struct {
		Description string
		Replay      string
	}{
		Description: report.Description,
		Replay:      report.Replay,
	}

	sendErr := SendEmailWithHTMLTemplate(userEmail, "Report Replay Details", "Controller/pages/ReportReplay.html", emailData)
	if sendErr != nil {
		return fmt.Errorf("error sending replay email: %v", sendErr)
	}

	log.Println("Replay email sent successfully!")
	return nil
}
