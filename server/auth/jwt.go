package auth

import (
	"crypto/rsa"
	"fmt"
	"log"

	jwt "github.com/dgrijalva/jwt-go"
	"golang.org/x/net/context"
	"google.golang.org/grpc/metadata"
)

// GetTokenFromContext return one jwt.Token from Context or false if the token is not present
// or is invalid.
func GetTokenFromContext(ctx context.Context, jwtPublicKey *rsa.PublicKey) (*jwt.Token, bool) {
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
		return jwtPublicKey, nil
	})
	if err == nil && jwtToken.Valid {
		return jwtToken, true
	}
	return nil, false
}
