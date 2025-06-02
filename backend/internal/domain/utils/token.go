package utils

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTClaims struct {
	ID    string
	Email string
	Role  int
}

func CreateTokenFromClaims(c JWTClaims) (string, error) {
	exp, err := time.ParseDuration(os.Getenv("TOKEN_LIFETIME"))
	if err != nil {
		return "", err
	}

	claims := jwt.MapClaims{
		"id":    c.ID,
		"role":  c.Role,
		"email": c.Email,
		"exp":   time.Now().Add(exp).Unix(),
	}

	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256, claims,
	)

	return token.SignedString([]byte(os.Getenv("TOKEN_SECRET")))
}
