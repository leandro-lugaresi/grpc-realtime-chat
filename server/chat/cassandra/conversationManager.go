package cassandra

import (
	"github.com/gocql/gocql"
	pb "github.com/leandro-lugaresi/grpc-realtime-chat/server/chat/chatpb"
	"github.com/pkg/errors"
)

type ConversationManager struct {
	Session *gocql.Session
}

func (m ConversationManager) CreateSchema(schema string) error {
	err := m.Session.Query("CREATE KEYSPACE IF NOT EXISTS " + schema +
		" WITH REPLICATION = { 'class' : 'SimpleStrategy', 'replication_factor' : 1 }").Exec()
	if err != nil {
		return errors.Wrap(err, "Failed to create the KEYSPACE")
	}

	err = m.Session.Query("CREATE TABLE IF NOT EXISTS " + schema + ".conversations (" +
		"id uuid PRIMARY KEY," +
		"title varchar," +
		"type varint," +
		"members set<uuid>," +
		"image blob," +
		"created_at timestamp," +
		"PRIMARY KEY(username))").Exec()
	if err != nil {
		return errors.Wrap(err, "Failed creating table")
	}

	err = m.Session.Query("create index on " + schema + ".conversations (members)").Exec()

	return err
}

func (m ConversationManager) GetByUserID(userID string, limit int32, offset int32) ([]*pb.Conversation, error) {
	iter := m.Session.Query(`SELECT id, title, type, members, image, created_at FROM conversations WHERE members CONTAINS ?`, userID).Iter()
	cl := []*pb.Conversation{}
	c := pb.Conversation{}
	for iter.Scan(&c.Id, &c.Title, &c.Type, &c.MemberIds, &c.Image, &c.CreationTime) {
		cl = append(cl, &c)
	}
	if err := iter.Close(); err != nil {
		return nil, err
	}
	return cl, nil
}

func (m ConversationManager) GetByID(id string) (*pb.Conversation, error) {
	c := &pb.Conversation{}
	q := m.Session.Query(`SELECT id, title, type, members, image, created_at FROM conversations WHERE id = ? LIMIT 1`, id)
	q.Consistency(gocql.One)
	err := q.Scan(c.Id, c.Title, c.Type, c.MemberIds, c.Image, c.CreationTime)
	return c, err
}

func (m ConversationManager) Create(c *pb.Conversation) error {
	id, err := gocql.RandomUUID()
	if err != nil {
		id = gocql.TimeUUID()
	}

	return m.Session.Query("INSERT INTO conversations (id, title, type, members, image, created_at) VALUES (?, ?, ?, ?, ?, ?)", id, c.Title, c.Type, c.MemberIds, c.Image).Exec()
}

func (m ConversationManager) RemoveMember(cID, memberID string) error {
	return m.Session.Query("UPDATE conversations SET members = members - {?} where id = ?", memberID, cID).Exec()
}

func (m ConversationManager) AddMember(cID, memberID string) error {
	return m.Session.Query("UPDATE conversations SET members = members + {?} where id = ?", memberID, cID).Exec()
}

func (m ConversationManager) Update(c *pb.Conversation) error {
	return m.Session.Query("UPDATE conversations SET title = ?, image = ? where id = ?", c.Title, c.Image, c.Id).Exec()
}
