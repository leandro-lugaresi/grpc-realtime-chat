package chat

import (
	"crypto/rsa"
	"io"
	"sync"

	jwt "github.com/dgrijalva/jwt-go"
	google_protobuf "github.com/golang/protobuf/ptypes/empty"
	"github.com/leandro-lugaresi/grpc-realtime-chat/server/auth"
	pb "github.com/leandro-lugaresi/grpc-realtime-chat/server/chat/chatpb"
	"github.com/nats-io/nats"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

type MessageManager interface {
	CreateSchema(schema string) error
	GetMessages(conversationID string, limit int32, offset int32) ([]*pb.Message, error)
	ReadMessages(cID string, mID string) error
	SaveMessage(ms *pb.Message, conversationID string) error
}

type MessageService struct {
	jwtPublicKey        *rsa.PublicKey
	messageManager      MessageManager
	conversationManager ConversationManager
	nc                  *nats.EncodedConn
}

func NewMessageService(rsaPublicKey []byte, m MessageManager, c ConversationManager, nc *nats.EncodedConn) (*MessageService, error) {
	pk, err := jwt.ParseRSAPublicKeyFromPEM(rsaPublicKey)
	if err != nil {
		return nil, errors.Wrap(err, "Error parsing the jwt public key")
	}
	return &MessageService{pk, m, c, nc}, nil
}

// GetHistory return the history of messages for an conversation in DESC order.
func (s *MessageService) GetHistory(cx context.Context, r *pb.GetHistoryRequest) (*pb.GetHistoryResponse, error) {
	m, err := s.messageManager.GetMessages(r.ConversationId, r.Limit, r.Offset)
	if err != nil {
		return nil, grpc.Errorf(codes.Internal, err.Error())
	}
	return &pb.GetHistoryResponse{Messages: m}, nil
}

// ReadHistory notifies the reading of messages from a conversation.
func (s *MessageService) ReadHistory(cx context.Context, r *pb.ReadHistoryRequest) (*google_protobuf.Empty, error) {
	return nil, grpc.Errorf(codes.Unimplemented, "Method not implemented")
}

// Comunicate Send and Receive Messages or events to/from conversations.
func (s *MessageService) Comunicate(str pb.MessageService_ComunicateServer) error {
	ctx := str.Context()
	id, err := auth.GetUserIDAuthenticated(ctx, s.jwtPublicKey)
	if err != nil {
		return err
	}
	msRecv := make(chan *pb.ChatMessage, 10)
	msToSend := make(chan *pb.ChatMessage, 10)
	errC := make(chan error, 1)
	closing := make(chan bool)
	s.nc.BindRecvChan("message.*."+id, msToSend)

	go func() {
		for {
			select {
			case <-closing:
				return
			default:
				ms, err := str.Recv()
				if err == io.EOF {
					close(msRecv)
					return
				}
				if err != nil {
					errC <- err
					close(msRecv)
					return
				}
				msRecv <- ms
			}
		}
	}()

	for {
		select {
		case ms := <-msToSend:
			str.Send(ms)
			continue
		case ms := <-msRecv:
			if ms.GetMessage() != nil {
				s.messageManager.SaveMessage(ms.GetMessage(), ms.ConversationId)
			}
			conv, err := s.conversationManager.GetByID(ms.ConversationId)
			if err != nil {
				continue
			}
			var wg sync.WaitGroup
			for _, mID := range conv.MemberIds {
				wg.Add(1)
				go func(mID string, id string, ms *pb.ChatMessage) {
					defer wg.Done()
					if mID == id {
						return
					}
					s.nc.Publish("message."+ms.ConversationId+"."+mID, ms)
				}(mID, id, ms)
			}
			wg.Wait()
			s.nc.Flush()
			continue
		case err := <-errC:
			return err
		case <-ctx.Done():
			close(msToSend)
			closing <- true
			return ctx.Err()
		}
	}
}
