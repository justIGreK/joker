package main

import (
	"jokegen/cmd/subscriber/consumer"
	"jokegen/internal/service"
	"jokegen/internal/storage"
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/nats-io/nats.go"
)

func main() {
	db, err := sqlx.Open("postgres", "postgresql://postgres:qwerty@localhost:5432/postgres?sslmode=disable")
	if err != nil {
		log.Fatal("error during opening db")
	}
	repos := storage.NewUsersPostgres(db)
	service := service.NewUsersService(repos)
	subscriber := consumer.NewPublish(service)

	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Println("nats is not connected")
	}
	
	subscriber.SubscribeNats(nc)

	select {}

}
