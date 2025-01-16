package service

import (
	"context"
	"log"

	"github.com/mohamedfawas/rmshop-auth-service/internal/config"
	"github.com/mohamedfawas/rmshop-auth-service/internal/domain"
	"github.com/mohamedfawas/rmshop-auth-service/internal/util"
	authv1 "github.com/mohamedfawas/rmshop-proto/gen/v1/auth"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type authService struct {
	repo      domain.AuthRepository
	jwtSecret string
	authv1.UnimplementedAuthServiceServer
}

func NewAuthService(repo domain.AuthRepository, config *config.Config) authv1.AuthServiceServer {
	service := &authService{
		repo:      repo,
		jwtSecret: config.JWT.Secret,
	}

	// Initialize admin user
	if err := service.initializeAdmin(); err != nil {
		log.Printf("Failed to initialize admin: %v", err)
	}

	return service
}

func (s *authService) initializeAdmin() error {
	passwordHash, err := util.HashPassword("admin123")
	if err != nil {
		return err
	}

	return s.repo.InitializeAdmin(context.Background(), "admin@rmshop.com", passwordHash)
}

func (s *authService) Login(ctx context.Context, req *authv1.LoginRequest) (*authv1.LoginResponse, error) {
	user, err := s.repo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, status.Error(codes.NotFound, "user not found")
	}

	if err := util.ComparePasswords(user.PasswordHash, req.Password); err != nil {
		return nil, status.Error(codes.Unauthenticated, "invalid credentials")
	}

	token, err := util.GenerateToken(user.ID, user.UserType, s.jwtSecret)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to generate token")
	}

	return &authv1.LoginResponse{
		Token:    token,
		UserId:   user.ID,
		UserType: user.UserType,
	}, nil
}

func (s *authService) Logout(ctx context.Context, req *authv1.LogoutRequest) (*authv1.LogoutResponse, error) {
	if err := s.repo.BlacklistToken(ctx, req.Token); err != nil {
		return nil, status.Error(codes.Internal, "failed to logout")
	}

	return &authv1.LogoutResponse{Success: true}, nil
}

func (s *authService) ValidateToken(ctx context.Context, req *authv1.ValidateTokenRequest) (*authv1.ValidateTokenResponse, error) {
	// Check if token is blacklisted
	isBlacklisted, err := s.repo.IsTokenBlacklisted(ctx, req.Token)
	if err != nil || isBlacklisted {
		return &authv1.ValidateTokenResponse{IsValid: false}, nil
	}

	// Validate token
	claims, err := util.ValidateToken(req.Token, s.jwtSecret)
	if err != nil {
		return &authv1.ValidateTokenResponse{IsValid: false}, nil
	}

	return &authv1.ValidateTokenResponse{
		IsValid:  true,
		UserId:   claims.UserID,
		UserType: claims.UserType,
	}, nil
}
