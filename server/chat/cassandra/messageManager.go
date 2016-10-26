package cassandra

import (
	"time"

	"github.com/gocql/gocql"
	"github.com/golang/protobuf/ptypes/timestamp"
	pb "github.com/leandro-lugaresi/grpc-realtime-chat/server/chat/chatpb"
	"github.com/pkg/errors"
)

type MessageManager struct {
	Session *gocql.Session
}

func (m MessageManager) CreateSchema(schema string) error {
	err := m.Session.Query("CREATE KEYSPACE IF NOT EXISTS " + schema +
		" WITH REPLICATION = { 'class' : 'SimpleStrategy', 'replication_factor' : 1 }").Exec()
	if err != nil {
		return errors.Wrap(err, "Failed to create the KEYSPACE")
	}

	err = m.Session.Query("CREATE TABLE IF NOT EXISTS " + schema + ".messages (" +
		"id uuid," +
		"sender_id uuid," +
		"conversation_id uuid," +
		"text text," +
		"audio blob," +
		"image blob," +
		"event int," +
		"extra_params map<string,string>," +
		"readed map<uuid,timestamp>," +
		"created_at timestamp," +
		"delivery_at timestamp," +
		"PRIMARY KEY  ((id), conversation_id ) )").Exec()
	if err != nil {
		return errors.Wrap(err, "Failed creating table")
	}

	err = m.Session.Query("create index on " + schema + ".messages (conversation_id)").Exec()

	return err
}

func (m MessageManager) GetMessages(conversationID string, limit int32, offset int32) ([]*pb.Message, error) {
	iter := m.Session.Query(`SELECT id, sender_id, text, audio, image, event, extra_params, created_at, delivery_at FROM messages WHERE conversation_id = ?`, conversationID).Iter()
	cl := []*pb.Message{}
	var (
		id          gocql.UUID
		sID         gocql.UUID
		text        string
		audio       []byte
		image       []byte
		event       int
		extraParams map[string]string
		created     time.Time
		delivery    time.Time
	)
	for iter.Scan(&id, &sID, &text, &audio, &image, &event, &extraParams, &created, &delivery) {
		ms := &pb.Message{}
		ms.Id = id.String()
		ms.SenderId = sID.String()
		ms.CreationTime = &timestamp.Timestamp{Seconds: int64(created.Second()), Nanos: int32(created.Nanosecond())}
		ms.DeliveryTime = &timestamp.Timestamp{Seconds: int64(delivery.Second()), Nanos: int32(delivery.Nanosecond())}

		if len(audio) > 0 {
			ms.Content = &pb.Message_Audio{Audio: audio}
		}
		if len(image) > 0 {
			ms.Content = &pb.Message_Image{Image: image}
		}
		if event >= 0 {
			ms.Content = &pb.Message_Event{Event: &pb.ActionEvent{Event: pb.ActionEvent_EventType(event)}}
		}
		cl = append(cl, ms)
	}
	if err := iter.Close(); err != nil {
		return nil, err
	}
	return cl, nil
}

func (m MessageManager) SaveMessage(message *pb.Message) error {
	return nil
}
