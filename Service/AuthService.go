package Service

import repository "main/Repository"

func ApplicationLogin(email string, password string) string {

	if repository.FindUserEmail(email) {
		if !repository.ValidUser(email, password) {
			return "Incorrect email or password"
		}
		return "Success"
	} else {
		return "Email not exist"
	}

}
