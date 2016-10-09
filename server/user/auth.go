package user

import (
	"crypto/rsa"
	"fmt"

	"time"

	jwt "github.com/dgrijalva/jwt-go"
	pb "github.com/leandro-lugaresi/grpc-realtime-chat/server/user/userpb"
	"github.com/pkg/errors"
	"github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

type AuthServer struct {
	JwtPrivateKey *rsa.PrivateKey
	UserService   UserService
}

type User struct {
	pb.User
	Password       string
	CreatedAt      int64
	UpdatedAt      int64
	LastActivityAt int64
}

type AuthClaims struct {
	ID string `json:"id"`
	jwt.StandardClaims
}

// UserService represets a service for User operations
type UserService interface {
	GetUserByUsername(username string) (*User, error)
	CreateUser(*User) error
}

func NewAuthServer(rsaPrivateKey []byte, s UserService) (*AuthServer, error) {
	publickey, err := jwt.ParseRSAPrivateKeyFromPEM(rsaPrivateKey)
	if err != nil {
		return nil, fmt.Errorf("Error parsing the jwt public key: %s", err)
	}

	return &AuthServer{publickey, s}, nil
}

func (as *AuthServer) SignUp(cx context.Context, r *pb.SignUpRequest) (*pb.Token, error) {
	pass, err := bcrypt.GenerateFromPassword([]byte(r.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to encrypt the Password")
	}
	t := time.Now().Unix()
	user := &User{
		User: pb.User{
			Id:       uuid.NewV4().String(),
			Name:     r.Name,
			Username: r.Username,
		},
		Password:       string(pass),
		CreatedAt:      t,
		UpdatedAt:      t,
		LastActivityAt: t,
	}
	err = as.UserService.CreateUser(user)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to signUp")
	}
	token, err := as.generateToken(user)
	if err != nil {
		return nil, err
	}
	return &pb.Token{token}, nil
}

func (as *AuthServer) SignIn(cx context.Context, r *pb.SignInRequest) (*pb.Token, error) {
	u, err := as.UserService.GetUserByUsername(r.Username)
	if err != nil {
		return nil, grpc.Errorf(codes.PermissionDenied, "Username or password invalid")
	}
	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(r.Password))
	if err != nil {
		return nil, grpc.Errorf(codes.PermissionDenied, "Username or password invalid")
	}

	token, err := as.generateToken(u)
	if err != nil {
		return nil, err
	}
	return &pb.Token{token}, nil
}

func (as *AuthServer) generateToken(user *User) (string, error) {
	claims := AuthClaims{
		user.Id,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 96).Unix(),
			Issuer:    "auth.service",
			IssuedAt:  time.Now().Unix(),
			Subject:   user.Username,
		},
	}

	t := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	ts, err := t.SignedString(as.JwtPrivateKey)
	if err != nil {
		return "", errors.Wrap(err, "Failed to create the auth token")
	}

	return ts, nil
}
