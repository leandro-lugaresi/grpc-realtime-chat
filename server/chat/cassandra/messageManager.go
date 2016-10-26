package cassandra

import (
	"log"
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
		log.Printf("received event: %v", event)
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

func (m MessageManager) SaveMessage(ms *pb.Message, conversationID string) error {
	id, err := gocql.RandomUUID()
	if err != nil {
		id = gocql.TimeUUID()
	}
	event := ms.GetEvent()
	e := -1
	if event != nil {
		e = int(event.Event)
	}
	cAt := time.Unix(ms.CreationTime.Seconds, int64(ms.CreationTime.Nanos))
	dAt := time.Unix(ms.DeliveryTime.Seconds, int64(ms.DeliveryTime.Nanos))
	return m.Session.Query("INSERT INTO messages "+
		"(id, sender_id, conversation_id, text, audio, image, event, extra_params, readed, created_at, delivery_at) "+
		"VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
		id, ms.SenderId, conversationID, ms.GetText(), ms.GetAudio(), ms.GetImage(), e, event.GetExtraParams(), nil, cAt, dAt).Exec()
}

func (m MessageManager) ReadMessages(cID string, mID string) error {
	return errors.New("Not implemented")
}
