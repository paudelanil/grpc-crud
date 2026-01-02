package main

import (
	"fmt"
	"log"
	"net"

	"github.com/paudelanil/grpc-crud/models"
	pb "github.com/paudelanil/grpc-crud/pb"
	"github.com/paudelanil/grpc-crud/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {

	dsn := "host=localhost user=postgres password=pass dbname=grpc_crud port=5432 sslmode=disable"

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})

	if err != nil {
		log.Fatal(err)
	}

	// Auto Migrate all tables at once
	if err := db.AutoMigrate(&models.Customer{}, &models.Account{}); err != nil {
		log.Fatalf("Failed to migrate: %v", err)
	}

	// start gRPC server
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%s", "localhost", "8090"))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()

	accountService := &service.AccountService{DB: db}

	pb.RegisterAccountServiceServer(grpcServer, accountService)

	reflection.Register(grpcServer)
	log.Println("gRPC server listening on port", "8090")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}
