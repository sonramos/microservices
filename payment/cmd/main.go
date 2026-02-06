package main

import (
	"log"

	"github.com/sonramos/microservices/payment/config"
	"github.com/sonramos/microservices/payment/internal/adapters/db"
	"github.com/sonramos/microservices/payment/internal/adapters/grpc"
	"github.com/sonramos/microservices/payment/internal/application/core/api"
)

func main() {
	dbAdapter, err := db.NewAdapter(config.GetDataSourceURL())
	if err != nil {
		log.Fatalf("Failed to connect to database. Error: %v", err)
	}
	application := api.NewApplication(dbAdapter)
	grpcAdapter := grpc.NewAdapter(application, config.GetApplicationPort())
	grpcAdapter.Run()
}
