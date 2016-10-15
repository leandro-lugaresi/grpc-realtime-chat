package mysql

import (
	"database/sql"

	"github.com/leandro-lugaresi/grpc-realtime-chat/server/user"
)

type UserManager struct {
	db *sql.DB
}

func (m UserManager) GetUserByUsername(username string) (*user.User, error) {
	u := &user.User{}
	err := m.db.QueryRow("SELECT id, name, username, password, created_at, updated_at last_activity_at FROM users WHERE username=?", username).Scan(
		u.Id,
		u.Name,
		u.Username,
		u.Password,
		u.CreatedAt,
		u.UpdatedAt,
		u.LastActivityAt)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (m UserManager) GetUserById(id string) (*user.User, error) {
	u := &user.User{}
	err := m.db.QueryRow("SELECT id, name, username, password, created_at, updated_at last_activity_at FROM users WHERE id=?", id).Scan(
		u.Id,
		u.Name,
		u.Username,
		u.Password,
		u.CreatedAt,
		u.UpdatedAt,
		u.LastActivityAt)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (m UserManager) UpdateUser(*user.User) error {

}

func (m UserManager) CreateUser(*user.User) error {
	stmp, err := m.db.Prepare("INSERT INTO users(`id`,`username`,`name`,`password`,`created_at`,`updated_at`,`last_activity_at`) VALUES (?, ?, ?, ?, ?, ?, ?)")
	if err != nil {

	}
}
func (m UserManager) FindUsersByUsernameOrName(username string, name string) ([]*user.User, error) {

}

func (m UserManager) FindUsersByIds(ids []string) ([]*user.User, error) {

}
