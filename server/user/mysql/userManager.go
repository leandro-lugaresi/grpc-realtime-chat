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

func (m UserManager) UpdateUser(u *user.User) error {
	stmp, err := m.db.Prepare("UPDATE users SET `username` = ?, `name` = ?, `password` = ?, `updated_at` = ?, `last_activity_at` = ? WHERE `id` = ?")
	if err != nil {
		return err
	}
	_, err = stmp.Exec(u.Username, u.Name, u.Password, u.UpdatedAt, u.LastActivityAt, u.Id)
	return err
}

func (m UserManager) CreateUser(u *user.User) error {
	stmp, err := m.db.Prepare("INSERT INTO users(`id`,`username`,`name`,`password`,`created_at`,`updated_at`,`last_activity_at`) VALUES (?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		return err
	}
	_, err = stmp.Exec(u.Id, u.Username, u.Name, u.Password, u.CreatedAt, u.UpdatedAt, u.LastActivityAt)
	return err
}

func (m UserManager) FindUsersByUsernameOrName(name string) ([]*user.User, error) {
	users := []*user.User{}
	rows, err := m.db.Query("SELECT id, name, username, password, created_at, updated_at last_activity_at FROM users WHERE username LIKE '?%' OR name LIKE '?%' LIMIT 100", name, name)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		u := &user.User{}
		err := rows.Scan(
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
		users = append(users, u)
	}

	if err != nil {
		return nil, err
	}
	return users, nil
}

func (m UserManager) FindUsersByIds(ids []string) ([]*user.User, error) {

}
