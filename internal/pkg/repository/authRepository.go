package repository

import (
	"context"
	"encoding/base64"
	"errors"
	"os"

	"github.com/golang-jwt/jwt/v5"
	"github.com/raulaguila/go-template/internal/pkg/domain"
)

func NewAuthRepository(userRepository domain.UserRepository) domain.AuthRepository {
	return &authRepository{
		userRepository: userRepository,
	}
}

type authRepository struct {
	userRepository domain.UserRepository
}

func (s *authRepository) Login(ctx context.Context, user *domain.User, ip string) (*domain.AuthResponse, error) {
	accessTime, refreshTime := "-", "-"
	if user.Expire {
		accessTime = os.Getenv("ACCESS_TOKEN_EXPIRE")
		refreshTime = os.Getenv("RFRESH_TOKEN_EXPIRE")
	}

	accessToken, err := user.GenerateToken(accessTime, os.Getenv("ACCESS_TOKEN_PRIVAT"), ip)
	if err != nil {
		return nil, err
	}

	refreshToken, err := user.GenerateToken(refreshTime, os.Getenv("RFRESH_TOKEN_PRIVAT"), ip)
	if err != nil {
		return nil, err
	}

	return &domain.AuthResponse{
		User: user,
		TokensResponse: domain.TokensResponse{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		},
	}, nil
}

func (s *authRepository) claims2user(ctx context.Context, parsedToken *jwt.Token, ip string) (*domain.User, error) {
	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok || !parsedToken.Valid {
		return nil, errors.New("invalid token")
	}

	usr, err := s.userRepository.GetUserByToken(ctx, claims["token"].(string))
	if err != nil {
		return nil, err
	}
	if val, ok := claims["ip"]; !ok || val.(string) != ip {
		return nil, domain.ErrInvalidIpAssociation
	}
	if val, ok := claims["expire"]; ok {
		usr.Expire = val.(bool)
	}

	if !usr.Status {
		return nil, errors.New("disabled user")
	}

	return usr, nil
}

func (s *authRepository) Me(ctx context.Context, userToken, base64Key, ip string) (*domain.User, error) {
	decodedKey, err := base64.StdEncoding.DecodeString(base64Key)
	if err != nil {
		return nil, err
	}

	parsedKey, err := jwt.ParseRSAPublicKeyFromPEM(decodedKey)
	if err != nil {
		return nil, err
	}

	parsedToken, err := jwt.Parse(userToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, err
		}

		return parsedKey, nil
	})
	if err != nil {
		return nil, err
	}

	return s.claims2user(ctx, parsedToken, ip)
}

func (s *authRepository) Refresh(ctx context.Context, user *domain.User, ip string) (*domain.TokensResponse, error) {
	accessTime, refreshTime := "-", "-"
	if user.Expire {
		accessTime = os.Getenv("ACCESS_TOKEN_EXPIRE")
		refreshTime = os.Getenv("RFRESH_TOKEN_EXPIRE")
	}

	accessToken, err := user.GenerateToken(accessTime, os.Getenv("ACCESS_TOKEN_PRIVAT"), ip)
	if err != nil {
		return nil, err
	}

	refreshToken, err := user.GenerateToken(refreshTime, os.Getenv("RFRESH_TOKEN_PRIVAT"), ip)
	if err != nil {
		return nil, err
	}

	return &domain.TokensResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *authRepository) GetUserByMail(ctx context.Context, userMail string) (*domain.User, error) {
	return s.userRepository.GetUserByMail(ctx, userMail)
}
