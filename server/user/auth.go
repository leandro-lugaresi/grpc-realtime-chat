package user

import (
	"crypto/rsa"
	"fmt"

	jwt "github.com/dgrijalva/jwt-go"
	pb "github.com/leandro-lugaresi/grpc-realtime-chat/server/user/userpb"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/net/context"
)

type authServer struct {
	jwtPrivateKey *rsa.PrivateKey
}

type User struct {
	pb.User
	Password       string
	CreatedAt      int64
	UpdatedAt      int64
	LastActivityAt int64
}

// UserService represets a service for User operations
type UserService interface {
	GetUserByUsername(username string) (*User, error)
	CreateUser(*User) error
}

func NewAuthServer(rsaPrivateKey []byte) (*authServer, error) {
	publickey, err := jwt.ParseRSAPrivateKeyFromPEM(rsaPrivateKey)
	if err != nil {
		return nil, fmt.Errorf("Error parsing the jwt public key: %s", err)
	}

	return &authServer{publickey}, nil
}

func (as *authServer) SignUp(cx context.Context, r *pb.SignUpRequest) (*pb.Token, error) {
	pass, err := bcrypt.GenerateFromPassword([]byte(r.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to encrypt the Password")
	}
	user := &User{
		User: pb.User{
			Name:     r.Name,
			Username: r.Username,
		},
		Password: string(pass),
	}

}

func (as *authServer) SignIn(cx context.Context, r *pb.SignInRequest) (*pb.Token, error) {

}
