package storage

import (
	"jokegen/internal/models"

	"github.com/jmoiron/sqlx"
)

type Users interface {
	CreateUser(login, password string) (int, error)
	GetUser(login, password string) (models.User, error)
	GetUserById(id int) (models.User, error)
	UpdateUserAttempts(id, attempts int) error
	UpdateUserAttemptsByLogin(login string, attempts int) error
}

type Store struct {
	Users
}

func NewStore(db *sqlx.DB) *Store {
	return &Store{
		Users: NewUsersPostgres(db),
	}
}
