package user_test

import (
	"io/ioutil"

	"testing"

	"github.com/leandro-lugaresi/grpc-realtime-chat/server/user"
)

type UserService struct {
	CreateUserFn       func(*user.User) error
	CreatedUserInvoked bool

	GetUserByUsernameFn      func(username string) (*user.User, error)
	GetUserByUsernameInvoked bool
}

func (s *UserService) CreateUser(user *user.User) error {
	s.CreatedUserInvoked = true
	return s.CreateUserFn(user)
}

func (s *UserService) GetUserByUsername(username string) (*user.User, error) {
	s.GetUserByUsernameInvoked = true
	return s.GetUserByUsernameFn(username)
}

func NewAuthServer() (*user.AuthServer, error) {
	keyData, _ := ioutil.ReadFile("../test/rsa_sample_key.pub")
	s := UserService
	return user.NewAuthServer(keyData, s)
}

func TestAuthService_SignUp(t *testing.T) {
	c, err := NewAuthServer()

	c.UserService.CreateUserFn = func(*user.User) error {
		return nil
	}
	c.Si
}
