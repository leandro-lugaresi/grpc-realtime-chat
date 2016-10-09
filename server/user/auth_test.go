package user_test

import (
	"io/ioutil"
	"net"

	"testing"
	"time"

	"github.com/leandro-lugaresi/grpc-realtime-chat/server/user"
	pb "github.com/leandro-lugaresi/grpc-realtime-chat/server/user/userpb"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type UserService struct {
	mock.Mock
}

func (s *UserService) CreateUser(user *user.User) error {
	ret := s.Mock.Called()
	return ret.Error(0)
}

func (s *UserService) GetUserByUsername(username string) (*user.User, error) {
	ret := s.Called()

	var r0 *user.User
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(*user.User)
	}
	r1 := ret.Error(1)

	return r0, r1
}

func TestServerAuthSuite(t *testing.T) {
	suite.Run(t, &ServerAuthTestSuite{})
}

type ServerAuthTestSuite struct {
	suite.Suite

	serverListener net.Listener
	server         *grpc.Server
	clientConn     *grpc.ClientConn
	testClient     pb.AuthServiceClient
	ctx            context.Context
	userService    UserService
}

func (s *ServerAuthTestSuite) SetupSuite() {
	var err error

	s.serverListener, err = net.Listen("tcp", "127.0.0.1:0")
	require.NoError(s.T(), err, "must be able to allocate a port for serverListener")
	//Create the service
	keyData, err := ioutil.ReadFile("../test/rsa_sample_key")
	require.NoError(s.T(), err, "must be able to read the rsa_key for tests")
	s.userService = UserService{}
	serv, err := user.NewAuthServer(keyData, &s.userService)
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

func (s *ServerAuthTestSuite) TearDownSuite() {
	if s.serverListener != nil {
		s.server.Stop()
		s.T().Logf("stopped grpc.Server at: %v", s.serverListener.Addr().String())
		s.serverListener.Close()

	}
	if s.clientConn != nil {
		s.clientConn.Close()
	}
}

func (s *ServerAuthTestSuite) SetupTest() {
	s.ctx, _ = context.WithTimeout(context.TODO(), 1*time.Second)
}

func (s *ServerAuthTestSuite) TestSignUp() {
	s.userService.On("CreateUser").Return(nil)
	r := pb.SignUpRequest{
		Name:     "Jhon Doe",
		Username: "jhon_doe",
		Password: "foo-123",
	}
	token, err := s.testClient.SignUp(s.ctx, &r)
	require.NoError(s.T(), err, "must not error on SignUp")
	assert.NotEmpty(s.T(), token, "Token must have a value")
}

func (s *ServerAuthTestSuite) TestSignIn() {
	s.userService.On("GetUserByUsername").Return(&user.User{
		User: pb.User{
			Id:       "118f3f16-84ed-4a4a-923f-77e4ffde04b6",
			Name:     "Jhon Doe",
			Username: "jhon_doe",
		},
		Password: "$2a$10$vmME6Fhw26sbyuYmzWKQIulu8dgPoH7KGmZVTpvCZqHiU7DE33X5S",
	}, nil)
	r := pb.SignInRequest{
		Username: "jhon_doe",
		Password: "foo-123",
	}
	token, err := s.testClient.SignIn(s.ctx, &r)
	require.NoError(s.T(), err, "must not error on SignUp")
	assert.NotEmpty(s.T(), token, "Token must have a value")
}

func (s *ServerAuthTestSuite) TestSignInShouldReturnErrorWhenPasswordsAreNotEqual() {
	s.userService.On("GetUserByUsername").Return(&user.User{
		User: pb.User{
			Id:       "118f3f16-84ed-4a4a-923f-77e4ffde04b6",
			Name:     "Jhon Doe",
			Username: "jhon_doe",
		},
		Password: "$2a$10$vmME6Fhw26sbyuYmzWKQIulu8dgPoH7KGmZVTpvCZqHiU7DE33X5S",
	})
	r := pb.SignInRequest{
		Username: "jhon_doe",
		Password: "foo-1234",
	}
	token, err := s.testClient.SignIn(s.ctx, &r)
	require.Error(s.T(), err, "Must return an error")
	assert.Empty(s.T(), token, "Must return an Empty token")
}
