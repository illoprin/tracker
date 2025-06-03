package utils

import (
	"os"
	"time"
	"tracker-backend/internal/config"

	"github.com/golang-jwt/jwt/v5"
)

type JWTClaims struct {
	UserID    string
	SessionID string
	Role      int
}

func CreateTokenFromClaims(c JWTClaims) (string, error) {
	exp, err := time.ParseDuration(os.Getenv(config.TokenLifetimeEnvName))
	if err != nil {
		return "", err
	}

	claims := jwt.MapClaims{
		"id":      c.UserID,
		"role":    c.Role,
		"session": c.SessionID,
		"exp":     time.Now().Add(exp).Unix(),
	}

	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256, claims,
	)

	return token.SignedString([]byte(os.Getenv(config.TokenSecretEnvName)))
}

func DecodeToken(token string) (*jwt.Token, *JWTClaims, error) {
	claims := jwt.MapClaims{}
	decoded, _, err := jwt.NewParser().ParseUnverified(token, claims)
	if err != nil {
		return nil, nil, err
	}

	claimsStruct := JWTClaims{
		UserID:    claims["id"].(string),
		SessionID: claims["session"].(string),
		Role:      int(claims["role"].(float64)),
	}

	return decoded, &claimsStruct, nil
}
