package service

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"jokegen/internal/models"
	"jokegen/internal/storage"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Users interface {
	CreateUser(login, password string) (int, error)
	GetUser(login, password string) (models.User, error)
	GetUserById(id int) (models.User, error)
	UpdateUserAttempts(id, attempts int) error
	UpdateUserAttemptsByLogin(login string, attempts int) error
}
type UsersService struct {
	User Users
}

type JokeResponse struct {
	Joke              Joke `json:"joke"`
	RemainingAttempts int  `json:"remainingAttempts"`
}
type Joke struct {
	Setup     string `json:"setup"`
	Punchline string `json:"punchline"`
}

const url string = "https://official-joke-api.appspot.com/random_joke"

func NewUsersService(users *storage.UsersPostgres) *UsersService {
	return &UsersService{User: users}
}

func (s *UsersService) CreateUser(login string, password string) error {
	_, err := s.User.CreateUser(login, password)
	if err != nil {
		return fmt.Errorf("error during creating acc: %w", err)
	}
	return nil
}

func (s *UsersService) LoginUser(login string, password string) (int, error) {
	user, err := s.User.GetUser(login, password)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, err
		}
		return 0, fmt.Errorf("error during getting acc: %w", err)
	}

	return user.Id, nil
}

func (s *UsersService) GetRandomJoke(userID int) (JokeResponse, error) {
	var response JokeResponse
	var joke Joke
	user, err := s.User.GetUserById(userID)
	if err != nil {
		return response, fmt.Errorf("error during getting userinfo: %w", err)
	}
	if user.Attempts == 0 {
		return response, errors.New("you have run out of attempts")
	}

	resp, err := http.Get(url)
	if err != nil {
		return response, fmt.Errorf("error during getting response: %w", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return response, fmt.Errorf("error during reading body: %w", err)
	}

	if err := json.Unmarshal(body, &joke); err != nil {
		return response, fmt.Errorf("error during unmarshal body: %w", err)
	}

	err = s.User.UpdateUserAttempts(userID, -1)
	if err != nil {
		return response, fmt.Errorf("error during updating attempts: %w", err)
	}
	response.Joke = joke
	response.RemainingAttempts = user.Attempts - 1
	return response, nil
}

func (s *UsersService) AddAttempts(userID, count int) error {
	err := s.User.UpdateUserAttempts(userID, count)
	if err != nil {
		return fmt.Errorf("error during adding attempts: %w", err)
	}
	return nil
}

func (s *UsersService) AddAttemptsByLogin(login string, count int) error {
	err := s.User.UpdateUserAttemptsByLogin(login, count)
	if err != nil {
		return fmt.Errorf("error during adding attempts: %w", err)
	}
	return nil
}

type tokenClaims struct {
	jwt.RegisteredClaims
	UserId int `json:"user_id"`
}

const (
	tokenTTL   = 12 * time.Hour
	signingKey = "fsjklj235OIUJlknm24"
)

func (s *UsersService) GenerateToken(login, password string) (string, error) {
	user, err := s.User.GetUser(login, password)
	if err != nil {
		return "", fmt.Errorf("error during getting user:%w", err)
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		}, user.Id,
	})

	entryToken, err := token.SignedString([]byte(signingKey))
	if err != nil {
		return "", fmt.Errorf("error during generating token: %w", err)
	}
	return entryToken, nil
}

func (s *UsersService) ParseToken(accessToken string) (int, error) {
	token, err := jwt.ParseWithClaims(accessToken, &tokenClaims{}, func(accessToken *jwt.Token) (interface{}, error) {
		if _, ok := accessToken.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return []byte(signingKey), nil
	})
	if err != nil {
		return 0, fmt.Errorf("error during parsing token: %w", err)
	}

	claims, ok := token.Claims.(*tokenClaims)
	if !ok {
		return 0, errors.New("token claims are not of type *tokenClaims")
	}

	return claims.UserId, nil
}
