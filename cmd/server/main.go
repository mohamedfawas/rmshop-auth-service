package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/lib/pq"
	"github.com/mohamedfawas/rmshop-auth-service/internal/config"
	"github.com/mohamedfawas/rmshop-auth-service/internal/repository"
	"github.com/mohamedfawas/rmshop-auth-service/internal/service"
	authv1 "github.com/mohamedfawas/rmshop-proto/gen/v1/auth"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize database connection
	db, err := sql.Open("postgres", fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.DBName,
		cfg.Database.SSLMode,
	))
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Test database connection
	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}
	log.Printf("Successfully connected to database")

	// Initialize repositories
	authRepo := repository.NewAuthRepository(db)

	// Initialize default admin
	adminEmail := "admin@example.com"
	adminPassword := "admin123"
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(adminPassword), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("Failed to hash admin password: %v", err)
	}

	if err := authRepo.InitializeAdmin(context.Background(), adminEmail, string(hashedPassword)); err != nil {
		log.Fatalf("Failed to initialize admin: %v", err)
	}
	log.Printf("Admin initialized with email: %s", adminEmail)

	// Initialize services
	authService := service.NewAuthService(authRepo, cfg)

	// Initialize gRPC server
	grpcServer := grpc.NewServer()
	authv1.RegisterAuthServiceServer(grpcServer, authService)

	// Enable reflection
	reflection.Register(grpcServer)

	// Start gRPC server
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	// Handle shutdown gracefully
	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
		<-sigCh
		log.Println("Received shutdown signal")
		grpcServer.GracefulStop()
	}()

	log.Printf("Starting auth service on %s:%s", cfg.Server.Host, cfg.Server.Port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
