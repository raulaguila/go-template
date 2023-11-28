package domain

import (
	"context"

	"github.com/raulaguila/go-template/pkg/validator"
)

type (
	AuthResponse struct {
		User         *User  `json:"user"`
		AccessToken  string `json:"accesstoken" validate:"jwt"`
		RefreshToken string `json:"refreshtoken" validate:"jwt"`
	}

	AuthRepository interface {
		Login(context.Context, *User) (*AuthResponse, error)
		Me(context.Context, string, string) (*User, error)
		Refresh(context.Context, *User) (*AuthResponse, error)
		GetUserByMail(context.Context, string) (*User, error)
	}

	AuthService interface {
		Login(context.Context, *User) (*AuthResponse, error)
		Me(context.Context, string, string) (*User, error)
		Refresh(context.Context, *User) (*AuthResponse, error)
		GetUserByMail(context.Context, string) (*User, error)
	}
)

func (s *AuthResponse) Validate() error {
	return validator.StructValidator.Validate(s)
}
