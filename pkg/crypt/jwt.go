package crypt

import (
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/kitae0522/gommunity/internal/config"
	"github.com/kitae0522/gommunity/pkg/exception"
)

func NewToken(userRole, userID string, secretKey []byte) (string, error) {
	expiration := time.Duration(config.Envs.JWTExpirationInSeconds) * time.Second
	claims := jwt.MapClaims{
		"role":      userRole,
		"uuid":      userID,
		"expiredAt": time.Now().Add(expiration).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secretKey)
}

func ParseJWT(jwtToken string) (string, error) {
	token, err := jwt.Parse(jwtToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, exception.ErrUnexpectedSigningMethod
		}
		return []byte(config.Envs.JWTSecret), nil
	})

	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if email, ok := claims["uuid"].(string); ok {
			return email, nil
		}
	}
	return "", exception.ErrInvalidTokenClaims
}
