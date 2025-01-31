package main

import (
	"jokegen/cmd/joker/handler"
	"jokegen/internal"
	"jokegen/internal/service"
	"jokegen/internal/storage"
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// @title Joker
// @version 1.0
// @description API Server for JOKER Application

// @host localhost:8000
// @BasePath /

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {

	db, err := sqlx.Open("postgres", "postgresql://postgres:qwerty@localhost:5432/postgres?sslmode=disable")
	if err != nil {
		log.Fatal("error during opening db")
	}
	repos := storage.NewUsersPostgres(db)
	handler := handler.Handler{Users: service.NewUsersService(repos)} 
	srv := new(internal.Server)
	if err := srv.Run("8000", handler.InitRoutes()); err != nil {
		log.Fatal("error during start server:", err.Error())
	}
}
