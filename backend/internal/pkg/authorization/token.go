package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func CreateAuthToken(
	userRole int, userID, userEmail, secret string,
) (string, error) {
	claims := jwt.MapClaims{
		"id":    userID,
		"role":  userRole,
		"email": userEmail,

		// token expires in 30 days
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
	}

	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256, claims,
	)

	return token.SignedString([]byte(secret))
}
