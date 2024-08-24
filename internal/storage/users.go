package storage

import (
	"fmt"
	"jokegen/internal/models"
	"log"

	"github.com/jmoiron/sqlx"
)

type UsersPostgres struct {
	db *sqlx.DB
}

func NewUsersPostgres(db *sqlx.DB) *UsersPostgres {
	return &UsersPostgres{db: db}
}

func (a *UsersPostgres) CreateUser(login, password string) (int, error) {
	var id int
	query := "INSERT INTO users (login, password) values ($1, $2)"
	_, err := a.db.Exec(query, login, password)
	if err != nil {
		return 0, fmt.Errorf("error during inserting into db: %w", err)
	}
	log.Print("user is created")
	return id, nil
}

func (a *UsersPostgres) GetUser(login, password string) (models.User, error) {
	var user, empty models.User
	query := "SELECT * FROM users WHERE login=$1 AND password=$2"
	err := a.db.Get(&user, query, login, password)
	if err != nil {
		return empty, err
	}
	return user, nil
}

func (u *UsersPostgres) GetUserById(id int) (models.User, error) {
	var user, empty models.User
	query := "SELECT * FROM users WHERE id = $1"
	err := u.db.Get(&user, query, id)
	if err != nil {
		return empty, fmt.Errorf("error during getting user by id:%w", err)
	}
	return user, nil
}

func (u *UsersPostgres) UpdateUserAttempts(id, attempts int) error {
	query := "UPDATE users SET attempts=attempts+$1 WHERE id=$2"
	_, err := u.db.Exec(query, attempts, id)
	if err != nil {
		return fmt.Errorf("error during updating user attempts: %w", err)
	}
	return nil
}

func (u *UsersPostgres) UpdateUserAttemptsByLogin(login string, attempts int) error {
	query := "UPDATE users SET attempts=attempts+$1 WHERE login=$2"
	_, err := u.db.Exec(query, attempts, login)
	if err != nil {
		return fmt.Errorf("error during updating user attempts: %w", err)
	}
	return nil
}