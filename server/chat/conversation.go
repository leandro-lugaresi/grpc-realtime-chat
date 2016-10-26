package chat

import (
	"crypto/rsa"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	google_protobuf "github.com/golang/protobuf/ptypes/empty"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/leandro-lugaresi/grpc-realtime-chat/server/auth"
	pb "github.com/leandro-lugaresi/grpc-realtime-chat/server/chat/chatpb"
	"github.com/pkg/errors"
	"github.com/satori/go.uuid"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

type ConversationManager interface {
	CreateSchema(schema string) error
	GetByUserID(userID string, limit int32, offset int32) ([]*pb.Conversation, error)
	GetByID(id string) (*pb.Conversation, error)
	Create(*pb.Conversation) error
	Update(*pb.Conversation) error
	AddMember(cID, memberID string) error
	RemoveMember(cID, memberID string) error
}

type ConversationService struct {
	jwtPublicKey        *rsa.PublicKey
	conversationManager ConversationManager
}

func NewConversationService(rsaPublicKey []byte, m ConversationManager) (*ConversationService, error) {
	pk, err := jwt.ParseRSAPublicKeyFromPEM(rsaPublicKey)
	if err != nil {
		return nil, errors.Wrap(err, "Error parsing the jwt public key")
	}
	return &ConversationService{pk, m}, nil
}

func (s *ConversationService) Get(ctx context.Context, r *pb.GetConversationsRequest) (*pb.GetConversationsResponse, error) {
	id, err := auth.GetUserIDAuthenticated(ctx, s.jwtPublicKey)
	if err != nil {
		return nil, err
	}
	c, err := s.conversationManager.GetByUserID(id, r.Limit, r.Offset)
	if err != nil {
		return nil, grpc.Errorf(codes.Internal, err.Error())
	}
	return &pb.GetConversationsResponse{c}, nil
}

func (s *ConversationService) Create(ctx context.Context, r *pb.CreateConversationRequest) (*pb.CreateConversationResponse, error) {
	t := &timestamp.Timestamp{Seconds: time.Now().Unix()}
	c := &pb.Conversation{
		Id:           uuid.NewV4().String(),
		Title:        r.Title,
		MemberIds:    r.MemberIds,
		Type:         r.Type,
		CreationTime: t,
		UpdateTime:   t,
	}
	err := s.conversationManager.Create(c)
	if err != nil {
		return nil, grpc.Errorf(codes.Internal, err.Error())
	}
	return &pb.CreateConversationResponse{c}, nil
}

func (s *ConversationService) Leave(ctx context.Context, r *pb.LeaveConversationRequest) (*google_protobuf.Empty, error) {
	id, err := auth.GetUserIDAuthenticated(ctx, s.jwtPublicKey)
	if err != nil {
		return nil, err
	}
	if err = s.conversationManager.RemoveMember(r.ConversationId, id); err != nil {
		return nil, grpc.Errorf(codes.Internal, err.Error())
	}
	return &google_protobuf.Empty{}, nil
}

func (s *ConversationService) AddMember(ctx context.Context, r *pb.MemberRequest) (*google_protobuf.Empty, error) {
	_, ok := auth.GetTokenFromContext(ctx, s.jwtPublicKey)
	if !ok {
		return nil, grpc.Errorf(codes.Unauthenticated, "valid token required.")
	}
	if err := s.conversationManager.AddMember(r.ConversationId, r.UserId); err != nil {
		return nil, grpc.Errorf(codes.Internal, err.Error())
	}
	return &google_protobuf.Empty{}, nil
}

func (s *ConversationService) RemoveMember(ctx context.Context, r *pb.MemberRequest) (*google_protobuf.Empty, error) {
	_, ok := auth.GetTokenFromContext(ctx, s.jwtPublicKey)
	if !ok {
		return nil, grpc.Errorf(codes.Unauthenticated, "valid token required.")
	}
	if err := s.conversationManager.RemoveMember(r.ConversationId, r.UserId); err != nil {
		return nil, grpc.Errorf(codes.Internal, err.Error())
	}
	return &google_protobuf.Empty{}, nil
}
