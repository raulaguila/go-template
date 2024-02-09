package domain

import (
	"context"
	"errors"

	"github.com/raulaguila/go-template/pkg/validator"
)

var ErrInvalidIpAssociation error = errors.New("invalid ip source")

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
		Login(context.Context, *User, string) (*AuthResponse, error)
		Me(context.Context, string, string, string) (*User, error)
		Refresh(context.Context, *User, string) (*TokensResponse, error)
		GetUserByMail(context.Context, string) (*User, error)
	}

	AuthService interface {
		Login(context.Context, *User, string) (*AuthResponse, error)
		Me(context.Context, string, string, string) (*User, error)
		Refresh(context.Context, *User, string) (*TokensResponse, error)
		GetUserByMail(context.Context, string) (*User, error)
	}
)

func (s *TokensResponse) Validate() error {
	return validator.StructValidator.Validate(s)
}
