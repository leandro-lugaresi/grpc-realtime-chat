package mysql

import (
	"database/sql"

	"github.com/leandro-lugaresi/grpc-realtime-chat/server/user"
)

type UserManager struct {
	db *sql.DB
}

func (m UserManager) GetUserByUsername(username string) (*user.User, error) {
	user := &user.User{}

	rows, err := db.Query("SELECT id, name, username, password,  FROM users WHERE username=?", username)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(&rawUser); err != nil {
			return nil, err
		}
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	err = proto.Unmarshal(rawUser, user)
	if err != nil {
		return nil, err
	}
	return user, nil
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
