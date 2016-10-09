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
	testClient     pb.AuthServiceClient
	ctx            context.Context
	userManager    UserManager
}

func (s *UserServiceSuite) SetupSuite() {
	var err error

	s.serverListener, err = net.Listen("tcp", "127.0.0.1:0")
	require.NoError(s.T(), err, "must be able to allocate a port for serverListener")
	//Create the service
	s.userManager = UserManager{}
	serv, err := user.NewUserServer(&s.userManager)
	require.NoError(s.T(), err, "must be able to create a authServer")
	// This is the point where we hook up the interceptor
	s.server = grpc.NewServer()
	pb.RegisterAuthServiceServer(s.server, serv)

	go func() {
		s.server.Serve(s.serverListener)
	}()

	s.clientConn, err = grpc.Dial(s.serverListener.Addr().String(), grpc.WithInsecure(), grpc.WithBlock(), grpc.WithTimeout(2*time.Second))
	require.NoError(s.T(), err, "must not error on client Dial")
	s.testClient = pb.NewAuthServiceClient(s.clientConn)
}

func (s *UserServiceSuite) TearDownSuite() {
	if s.serverListener != nil {
		s.server.Stop()
		s.T().Logf("stopped grpc.Server at: %v", s.serverListener.Addr().String())
		s.serverListener.Close()

	}
	if s.clientConn != nil {
		s.clientConn.Close()
	}
}

func (s *UserServiceSuite) SetupTest() {
	s.ctx, _ = context.WithTimeout(context.TODO(), 1*time.Second)
}
