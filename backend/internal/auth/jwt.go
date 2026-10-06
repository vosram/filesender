package auth

import (
	"encoding/base64"
	"errors"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/vosram/filesender/backend/internal/database"
)

type CustomJWTClaims struct {
	Role       string `json:"role"`
	Membership string `json:"membership"`
	jwt.RegisteredClaims
}

const TokenIssuer = "Filesender"
const AccessTokenDuration = 10 * time.Minute

// Creates a signed access JWT with a default of 10 minute expiration.
func CreateAccessJWT(user database.User, jwtSecret string) (string, error) {
	signingKey, err := base64.StdEncoding.DecodeString(jwtSecret)
	if err != nil {
		log.Printf("CreateAccessJWT: failed to decode base64 jwt secret")
		return "", err
	}

	claims := CustomJWTClaims{
		user.Role,
		user.Membership,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(AccessTokenDuration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    TokenIssuer,
			Subject:   user.ID.String(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(signingKey)
}

func ValidateAccessJWT(tokenStr string, jwtSecret string) (CustomJWTClaims, error) {
	signingKey, err := base64.StdEncoding.DecodeString(jwtSecret)
	if err != nil {
		log.Printf("CreateAccessJWT: failed to decode base64 jwt secret")
		return CustomJWTClaims{}, err
	}
	token, err := jwt.ParseWithClaims(tokenStr, &CustomJWTClaims{}, func(token *jwt.Token) (any, error) {
		return signingKey, nil
	}, nil)
	if err != nil {
		log.Printf("ValidateAccessJWT: Invalid Token")
		return CustomJWTClaims{}, err
	}
	if claims, ok := token.Claims.(*CustomJWTClaims); ok {
		return *claims, nil
	}
	log.Printf("ValidateAccessJWT: Token does not have expected claims")
	return CustomJWTClaims{}, errors.New("Token does not have expected claims")
}
