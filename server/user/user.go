package user

import (
	google_protobuf "github.com/golang/protobuf/ptypes/empty"
	pb "github.com/leandro-lugaresi/grpc-realtime-chat/server/user/userpb"

	_ "golang.org/x/crypto/bcrypt"
	"golang.org/x/net/context"
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
	CreateUser(*User) error
}

type UserService struct {
	UserManager UserManager
}

func NewUserService(s UserManager) *UserService {
	return &UserService{s}
}

func (s *UserService) ChangePassword(ctx context.Context, r *pb.ChangePasswordRequest) (*google_protobuf.Empty, error) {
	return &google_protobuf.Empty{}, nil
}

func (s *UserService) UpdateProfile(ctx context.Context, r *pb.UpdateProfileRequest) (*google_protobuf.Empty, error) {
	return &google_protobuf.Empty{}, nil
}

func (s *UserService) GetUsers(ctx context.Context, r *pb.GetUsersRequest) (*pb.GetUsersResponse, error) {
	return &pb.GetUsersResponse{}, nil
}
