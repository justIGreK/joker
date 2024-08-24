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
	repos := storage.NewStore(db)
	if err != nil {
		log.Fatal("error during opening db")
	}
	services := service.NewPublish(repos)
//	services := service.NewHandlerService(repos)
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Println("nats is not connected")
	}
	consumer.SubscribeNats(nc, services)
	
	select {}

}
