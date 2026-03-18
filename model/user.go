package model

import "sunhost/config"

type User struct {
	ID       int    `json:"id"`
	FullName string `json:"full_name"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (u *User) TableName() string {
	return "users"
}

func (u *User) GetID() int {
	return u.ID
}
func (u *User) Create() error {
	_, err := config.DB.Exec(`
		INSERT INTO users (full_name, username, email, password)
		VALUES (?, ?, ?, ?)
	`, u.FullName, u.Username, u.Email, u.Password)
	return err
}

func (u *User) GetByUsername(username string) error {
	return config.DB.QueryRow(`
		SELECT id, full_name, username, email, password
		FROM users
		WHERE username = ?
	`, username).Scan(&u.ID, &u.FullName, &u.Username, &u.Email, &u.Password)
}

func (u *User) GetByEmail(email string) error {
	return config.DB.QueryRow(`
		SELECT id, full_name, username, email, password
		FROM users
		WHERE email = ?
	`, email).Scan(&u.ID, &u.FullName, &u.Username, &u.Email, &u.Password)
}

func (u *User) DeleteUserByID() error {
	_, err := config.DB.Exec(`
		DELETE FROM users WHERE id = ?
	`, u.ID)
	return err
}
