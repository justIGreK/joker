package service

import (
	"jokegen/internal/storage"
)

type Users interface {
	CreateUser(login string, password string) error
	LoginUser(login string, password string) (int, error)
	GetRandomJoke(userID int) (JokeResponse, error)
	GenerateToken(login, password string) (string, error)
	ParseToken(accessToken string) (int, error)
	AddAttempts(userID, count int) error
	AddAttemptsByLogin(login string, count int) error
}

type Publisher interface {
	AddAttempts(userID, count int) error
}

type Service struct {
	Users
}

type Publish struct{
	Publisher
}

func NewPublish(store *storage.Store) *Publish{
	return &Publish{}
}
func NewService(store *storage.Store) *Service {
	return &Service{
		Users: NewUsersService(store.Users),
	}
	
}


