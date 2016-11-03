package auth

import (
	"crypto/rsa"
	"fmt"
	"log"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"google.golang.org/grpc/metadata"
)

// GetTokenFromContext return one jwt.Token from Context or false if the token is not present
// or is invalid.
func GetTokenFromContext(ctx context.Context, jwtPK *rsa.PublicKey) (*jwt.Token, bool) {
	md, ok := metadata.FromContext(ctx)
	if !ok {
		return nil, false
	}

	token, ok := md["authorization"]
	if !ok {
		return nil, false
	}

	jwtToken, err := jwt.ParseWithClaims(token[0], &jwt.StandardClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
			log.Printf("Unexpected signing method: %v", t.Header["alg"])
			return nil, fmt.Errorf("invalid token")
		}
		return jwtPK, nil
	})
	if err == nil && jwtToken.Valid {
		return jwtToken, true
	}
	return nil, false
}

func GetUserIDAuthenticated(ctx context.Context, jwtPK *rsa.PublicKey) (string, error) {
	token, ok := GetTokenFromContext(ctx, jwtPK)
	if !ok {
		return "", errors.New("valid token required")
	}
	claims := token.Claims.(*jwt.StandardClaims)
	return claims.Audience, nil
}
