package handlers

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func newSignedToken(Issuer string, expirationHours int, id int, JWTSecret string) (signedToken string, err error) {
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.RegisteredClaims{
			Issuer:   Issuer,
			IssuedAt: jwt.NewNumericDate(time.Now().UTC()),
			ExpiresAt: jwt.NewNumericDate(
				time.Now().UTC().Add(time.Duration(expirationHours) * time.Hour),
			),
			Subject: fmt.Sprintf("%d", id),
		},
	)

	return token.SignedString([]byte(JWTSecret))
}
