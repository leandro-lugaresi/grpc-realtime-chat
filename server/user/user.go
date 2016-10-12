package user

import (
	google_protobuf "github.com/golang/protobuf/ptypes/empty"
	pb "github.com/leandro-lugaresi/grpc-realtime-chat/server/user/userpb"
	"crypto/rsa"
	"fmt"
	"log"

	_ "golang.org/x/crypto/bcrypt"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"github.com/pkg/errors"
)

type User struct {
	pb.User
	Password       string
	CreatedAt      int64
	UpdatedAt      int64
	LastActivityAt int64
}

// UserManager represets a service for User operations
type UserManager interface {
	GetUserByUsername(username string) (*User, error)
	GetUserById(id string) (*User, error)
	CreateUser(*User) error
}

type UserService struct {
	UserManager UserManager,
	jwtPublicKey *rsa.PublicKey,
}

func NewUserService(rsaPublicKey []byte, s UserManager) (*UserService, error) {
	publickey, err := jwt.ParseRSAPublicKeyFromPEM(rsaPublicKey)
	if err != nil {
		return nil, errors.Wrap(err, "Error parsing the jwt public key")
	}
	return &UserService{s, publicKey}
}

func (s *UserService) ChangePassword(ctx context.Context, r *pb.ChangePasswordRequest) (*google_protobuf.Empty, error) {
	md, ok := metadata.FromContext(ctx)
	if !ok {
		return nil, grpc.Errorf(codes.Unauthenticated, "valid token required.")
	}

	jwtToken, ok := md["authorization"]
	if !ok {
		return nil, grpc.Errorf(codes.Unauthenticated, "valid token required.")
	}
	token, err = validateToken(jwtToken[0], hs.jwtPublicKey)
	if err != nil {
		return nil, grpc.Errorf(codes.Unauthenticated, "valid token required.")
	}
	user := s.UserManager.GetUserById()
	return &google_protobuf.Empty{}, nil
}

func (s *UserService) UpdateProfile(ctx context.Context, r *pb.UpdateProfileRequest) (*google_protobuf.Empty, error) {
	return &google_protobuf.Empty{}, nil
}

func (s *UserService) GetUsers(ctx context.Context, r *pb.GetUsersRequest) (*pb.GetUsersResponse, error) {
	return &pb.GetUsersResponse{}, nil
}

func validateToken(token string, publicKey *rsa.PublicKey) (*jwt.Token, error) {
	jwtToken, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
			log.Printf("Unexpected signing method: %v", t.Header["alg"])
			return nil, fmt.Errorf("invalid token")
		}
		return publicKey, nil
	})
	if err == nil && jwtToken.Valid {
		return jwtToken, nil
	}
	return nil, err
}
