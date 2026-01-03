package main

import (
	"fmt"
	"log"
	"net"

	"github.com/paudelanil/grpc-crud/internal/handler"
	"github.com/paudelanil/grpc-crud/internal/repository"
	"github.com/paudelanil/grpc-crud/internal/service"
	"github.com/paudelanil/grpc-crud/models"
	pb "github.com/paudelanil/grpc-crud/pb"
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
	if err := db.AutoMigrate(&models.Customer{}, &models.Account{}, &models.User{}); err != nil {
		log.Fatalf("Failed to migrate: %v", err)
	}

	// Initialize Repositories
	userRepo := repository.NewUserRepository(db)
	customerRepo := repository.NewCustomerRepository(db)
	accountRepo := repository.NewAccountRepository(db)

	// Initialize Services
	jwtSecret := "your-secret-key-change-this-in-production" // TODO: Move to environment variable
	authService := service.NewAuthService(userRepo, jwtSecret)
	customerService := service.NewCustomerService(customerRepo)
	accountService := service.NewAccountService(accountRepo, customerRepo)

	// Initialize Handlers
	accountHandler := handler.NewAccountHandler(customerService, accountService)
	authHandler := handler.NewAuthHandler(authService)

	// start gRPC server
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%s", "localhost", "8090"))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()

	// Register gRPC services
	pb.RegisterAccountServiceServer(grpcServer, accountHandler)
	pb.RegisterLoginServiceServer(grpcServer, authHandler)

	reflection.Register(grpcServer)
	log.Println("gRPC server listening on port", "8090")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}
