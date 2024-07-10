package main

import (
	"context"
	"log"
	"notification_service/server"
	"notification_service/service"
	s "notification_service/setup"

	_ "github.com/lib/pq"
)

func main() {
	ctx := context.Background()
	db, err := s.Setup(ctx)
	if err != nil {
		log.Fatalf("could not configure db: %v", err)
		return
	}
	s := service.NewNotificationService(db.Logger, db)

	// Create server
	server.ServerSetup(s)
}
