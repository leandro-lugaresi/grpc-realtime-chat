package mysql

import (
	"github.com/boltdb/bolt"

	"github.com/leandro-lugaresi/grpc-realtime-chat/server/user"
)

type UserManager struct {
	db *bolt.DB
}

func (m UserManager) GetUserByUsername(username string) (*user.User, error) {

}
func (m UserManager) GetUserById(id string) (*user.User, error) {

}
func (m UserManager) UpdateUser(*user.User) error {

}
func (m UserManager) CreateUser(*user.User) error {

}
func (m UserManager) FindUsersByUsernameOrName(username string, name string) ([]*user.User, error) {

}

func (m UserManager) FindUsersByIds(ids []string) ([]*user.User, error) {

}
