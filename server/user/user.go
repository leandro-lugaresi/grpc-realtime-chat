package user

import (
	"crypto/rsa"

	jwt "github.com/dgrijalva/jwt-go"
	google_protobuf "github.com/golang/protobuf/ptypes/empty"
	"github.com/leandro-lugaresi/grpc-realtime-chat/server/auth"
	pb "github.com/leandro-lugaresi/grpc-realtime-chat/server/user/userpb"

	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
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
	UpdateUser(*User) error
	CreateUser(*User) error
	FindUsersByUsernameOrName(name string) ([]*User, error)
	FindUsersByIds(ids []string) ([]*User, error)
}

type UserService struct {
	UserManager  UserManager
	jwtPublicKey *rsa.PublicKey
}

func NewUserService(rsaPublicKey []byte, s UserManager) (*UserService, error) {
	pk, err := jwt.ParseRSAPublicKeyFromPEM(rsaPublicKey)
	if err != nil {
		return nil, errors.Wrap(err, "Error parsing the jwt public key")
	}
	return &UserService{s, pk}, nil
}

func (s *UserService) ChangePassword(ctx context.Context, r *pb.ChangePasswordRequest) (*google_protobuf.Empty, error) {
	user, err := s.getUserFromContext(ctx)
	if err != nil {
		return nil, err
	}
	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(r.OldPassword)) != nil {
		return nil, grpc.Errorf(codes.PermissionDenied, "Password invalid")
	}
	pass, err := bcrypt.GenerateFromPassword([]byte(r.NewPassword), bcrypt.DefaultCost)
	user.Password = string(pass)
	err = s.UserManager.UpdateUser(user)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to save the user data")
	}
	return &google_protobuf.Empty{}, nil
}

func (s *UserService) UpdateProfile(ctx context.Context, r *pb.UpdateProfileRequest) (*google_protobuf.Empty, error) {
	user, err := s.getUserFromContext(ctx)
	if err != nil {
		return nil, err
	}
	if len(r.Name) > 0 {
		user.Name = r.Name
	}
	if len(r.Username) > 0 {
		user.Username = r.Username
	}
	err = s.UserManager.UpdateUser(user)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to save the user data")
	}
	return &google_protobuf.Empty{}, nil
}

func (s *UserService) GetUsers(ctx context.Context, r *pb.GetUsersRequest) (*pb.GetUsersResponse, error) {
	_, err := s.getUserFromContext(ctx)
	if err != nil {
		return nil, err
	}

	var users []*User
	if len(r.Ids) > 0 {
		users, err = s.UserManager.FindUsersByIds(r.Ids)
	}
	if len(r.Name) > 3 || len(r.Username) > 3 {
		users, err = s.UserManager.FindUsersByUsernameOrName(r.Username, r.Name)
	}
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to get the users")
	}
	rU := make([]*pb.User, len(users))
	for i, user := range users {
		rU[i] = &user.User
	}
	return &pb.GetUsersResponse{rU}, nil
}

func (s *UserService) getUserFromContext(ctx context.Context) (*User, error) {
	token, ok := auth.GetTokenFromContext(ctx, s.jwtPublicKey)
	if !ok {
		return nil, grpc.Errorf(codes.Unauthenticated, "valid token required.")
	}
	claims := token.Claims.(*jwt.StandardClaims)
	user, err := s.UserManager.GetUserById(claims.Audience)
	if err != nil {
		return nil, grpc.Errorf(codes.FailedPrecondition, "User not found")
	}
	return user, nil
}
