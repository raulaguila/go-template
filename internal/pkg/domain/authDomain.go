package domain

import (
	"context"

	"github.com/raulaguila/go-template/pkg/validator"
)

type (
	TokensResponse struct {
		AccessToken  string `json:"accesstoken" validate:"jwt"`
		RefreshToken string `json:"refreshtoken" validate:"jwt"`
	}

	AuthResponse struct {
		User *User `json:"user"`
		TokensResponse
	}

	AuthRepository interface {
		Login(context.Context, *User) (*AuthResponse, error)
		Me(context.Context, string, string) (*User, error)
		Refresh(context.Context, *User) (*TokensResponse, error)
		GetUserByMail(context.Context, string) (*User, error)
	}

	AuthService interface {
		Login(context.Context, *User) (*AuthResponse, error)
		Me(context.Context, string, string) (*User, error)
		Refresh(context.Context, *User) (*TokensResponse, error)
		GetUserByMail(context.Context, string) (*User, error)
	}
)

func (s *TokensResponse) Validate() error {
	return validator.StructValidator.Validate(s)
}
