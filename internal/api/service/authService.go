package service

import (
	"context"

	"github.com/raulaguila/go-template/internal/pkg/domain"
)

func NewAuthService(r domain.AuthRepository) domain.AuthService {
	return &authService{
		authRepository: r,
	}
}

type authService struct {
	authRepository domain.AuthRepository
}

func (s *authService) Login(ctx context.Context, user *domain.User, ip string) (*domain.AuthResponse, error) {
	return s.authRepository.Login(ctx, user, ip)
}

func (s *authService) Me(ctx context.Context, userToken, base64Key string, ip string) (*domain.User, error) {
	return s.authRepository.Me(ctx, userToken, base64Key, ip)
}

func (s *authService) Refresh(ctx context.Context, user *domain.User, ip string) (*domain.TokensResponse, error) {
	return s.authRepository.Refresh(ctx, user, ip)
}

func (s *authService) GetUserByMail(ctx context.Context, email string) (*domain.User, error) {
	return s.authRepository.GetUserByMail(ctx, email)
}
