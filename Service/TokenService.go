package Service

import (
	"encoding/json"
	"fmt"
	model "main/Model"
	repository "main/Repository"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var jwtKey = []byte("tajna_lozinka")

func GenerateToken(user model.User, exp time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":        user.ID,
		"email":     user.Email,
		"firstName": user.Firstname,
		"lastName":  user.Lastname,
		"phone":     user.Phonenumber,
		"password":  user.Password,
		"verified":  user.Verified,
		"location":  user.Location,
		"compnay":   user.CompmnayID,
		"role":      user.Role,
		"penalties": user.Penalies,
		"exp":       time.Now().Add(exp).Unix(),
	})

	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func SplitTokenHeder(authHeader string) string {

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return ""
	}
	tokenString := parts[1]
	return tokenString
}

// Parses the JWT token string and returns the token object.
func ParseTokenString(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtKey, nil
	})
	if err != nil {
		return token, fmt.Errorf("failed to parse token: %v", err)
	}
	return token, nil
}

func GetUserFromToken(res http.ResponseWriter, req *http.Request, token string) (model.User, *jwt.Token) {
	token = SplitTokenHeder(token)
	tokenpointer, _ := ParseTokenString(token)
	if tokenpointer == nil {
		var user model.User
		return user, nil
	}
	claims, _ := tokenpointer.Claims.(jwt.MapClaims)
	userID, _ := primitive.ObjectIDFromHex(claims["id"].(string))
	user, _ := repository.FindUserById(userID, res)
	return user, tokenpointer
}

func TokenAppLoginLogic(res http.ResponseWriter, req *http.Request, email string, password string) {
	message := ApplicationLogin(email, password)

	if message == "Success" {
		user, _ := repository.GetUserData(email)
		token, _ := GenerateToken(user, time.Hour)
		if user.Verified {
			res.WriteHeader(http.StatusOK)
			json.NewEncoder(res).Encode(token)
		} else {
			res.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(res).Encode("You must verify email before use app")
		}

	} else {
		res.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(res).Encode(message)
	}

}
