package chat

import (
	"crypto/rsa"

	jwt "github.com/dgrijalva/jwt-go"
	google_protobuf "github.com/golang/protobuf/ptypes/empty"
	pb "github.com/leandro-lugaresi/grpc-realtime-chat/server/chat/chatpb"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
)

type MessageManager interface {
	CreateSchema(schema string) error
	GetMessages(conversationID string, limit int32, offset int32) ([]*pb.Message, error)
	ReadMessages(cID string, mID string) error
	SaveMessage(*pb.Message) error
}

type MessageService struct {
	jwtPublicKey   *rsa.PublicKey
	messageManager MessageManager
}

func NewMessageService(rsaPublicKey []byte, m MessageManager) (*MessageService, error) {
	pk, err := jwt.ParseRSAPublicKeyFromPEM(rsaPublicKey)
	if err != nil {
		return nil, errors.Wrap(err, "Error parsing the jwt public key")
	}
	return &MessageService{pk, m}, nil
}

// Return the history of messages for an conversation in DESC order.
func (s *MessageService) GetHistory(cx context.Context, r *pb.GetHistoryRequest) (*pb.GetHistoryResponse, error) {
	return nil, nil
}

// Notifies the reading of messages from a channel or a user.
func (s *MessageService) ReadHistory(cx context.Context, r *pb.ReadHistoryRequest) (*google_protobuf.Empty, error) {
	return nil, nil
}

// Send and receive Messages or events to/from conversations.
func (s *MessageService) Comunicate(pb.MessageService_ComunicateServer) error {
	return nil
}
