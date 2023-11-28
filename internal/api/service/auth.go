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

func (s *authService) Login(ctx context.Context, user *domain.User) (*domain.AuthResponse, error) {
	return s.authRepository.Login(ctx, user)
}

func (s *authService) Me(ctx context.Context, userToken, base64Key string) (*domain.User, error) {
	return s.authRepository.Me(ctx, userToken, base64Key)
}

func (s *authService) Refresh(ctx context.Context, user *domain.User) (*domain.AuthResponse, error) {
	return s.authRepository.Refresh(ctx, user)
}

func (s *authService) GetUserByMail(ctx context.Context, email string) (*domain.User, error) {
	return s.authRepository.GetUserByMail(ctx, email)
}
