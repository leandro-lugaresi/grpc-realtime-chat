package user_test

import (
	"net"

	"testing"
	"time"

	"github.com/leandro-lugaresi/grpc-realtime-chat/server/user"
	pb "github.com/leandro-lugaresi/grpc-realtime-chat/server/user/userpb"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type UserManager struct {
	mock.Mock
}

func (s *UserManager) CreateUser(user *user.User) error {
	ret := s.Mock.Called()
	return ret.Error(0)
}

func (s *UserManager) GetUserByUsername(username string) (*user.User, error) {
	ret := s.Called()

	var r0 *user.User
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(*user.User)
	}
	r1 := ret.Error(1)

	return r0, r1
}

func TestUserServiceSuite(t *testing.T) {
	suite.Run(t, &UserServiceSuite{})
}

type UserServiceSuite struct {
	suite.Suite

	serverListener net.Listener
	server         *grpc.Server
	clientConn     *grpc.ClientConn
	testClient     pb.UserServiceClient
	ctx            context.Context
	userManager    UserManager
	jwtCreds       credentials.PerRPCCredentials
}

func (s *UserServiceSuite) SetupTest() {
	var err error

	s.serverListener, err = net.Listen("tcp", "127.0.0.1:0")
	require.NoError(s.T(), err, "must be able to allocate a port for serverListener")
	//Create the service
	s.userManager = UserManager{}
	serv := user.NewUserService(&s.userManager)
	// This is the point where we hook up the interceptor
	s.server = grpc.NewServer()
	pb.RegisterUserServiceServer(s.server, serv)

	go func() {
		s.server.Serve(s.serverListener)
	}()

	s.clientConn, err = grpc.Dial(
		s.serverListener.Addr().String(),
		grpc.WithInsecure(),
		grpc.WithPerRPCCredentials(),
		grpc.WithTimeout(2*time.Second)
	)
	require.NoError(s.T(), err, "must not error on client Dial")
	s.testClient = pb.NewUserServiceClient(s.clientConn)
	s.ctx, _ = context.WithTimeout(context.TODO(), 1*time.Second)
}

func (s *UserServiceSuite) TearDownTest() {
	if s.serverListener != nil {
		s.server.Stop()
		s.T().Logf("stopped grpc.Server at: %v", s.serverListener.Addr().String())
		s.serverListener.Close()

	}
	if s.clientConn != nil {
		s.clientConn.Close()
	}
}

func (s *UserServiceSuite) ChangePassword() {
	r := &pb.ChangePasswordRequest{
		NewPassword: "foooo-bar-baz",
		OldPassword: "foo-123",
	}
	s.ctx.
	s.testClient.ChangePassword(s.ctx, r)
}
