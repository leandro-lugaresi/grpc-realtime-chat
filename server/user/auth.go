package user

import (
	"crypto/rsa"
	"fmt"

	jwt "github.com/dgrijalva/jwt-go"
	pb "github.com/leandro-lugaresi/grpc-realtime-chat/server/user/userpb"
	"golang.org/x/net/context"
)

// UserService represets a service for User operations
type UserService interface {
	GetUserByUsername(username string) (*pb.User, error)
	CreateUser(*pb.User) error
}

type authServer struct {
	jwtPrivateKey *rsa.PrivateKey
}

func NewAuthServer(rsaPrivateKey []byte) (*authServer, error) {
	publickey, err := jwt.ParseRSAPrivateKeyFromPEM(rsaPrivateKey)
	if err != nil {
		return nil, fmt.Errorf("Error parsing the jwt public key: %s", err)
	}

	return &authServer{publickey}, nil
}

func (as *authServer) SignUp(cx context.Context, r *pb.SignUpRequest) (*pb.Token, error) {
	user := &pb.User{
		Username: r.Username,
		Name:     r.Name,
	}

}

func (as *authServer) SignIn(cx context.Context, r *pb.SignInRequest) (*pb.Token, error) {

}
