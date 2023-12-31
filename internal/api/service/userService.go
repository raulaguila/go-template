package service

import (
	"context"

	"github.com/raulaguila/go-template/internal/pkg/domain"
	"github.com/raulaguila/go-template/internal/pkg/dto"
	"github.com/raulaguila/go-template/pkg/filter"
)

func NewUserService(r domain.UserRepository) domain.UserService {
	return &userService{
		userRepository: r,
	}
}

type userService struct {
	userRepository domain.UserRepository
}

// Implementation of 'GetUsers'.
func (s *userService) GetUsersOutputDTO(ctx context.Context, filter *filter.UserFilter) (*dto.ItemsOutputDTO, error) {
	return s.userRepository.GetUsersOutputDTO(ctx, filter)
}

// Implementation of 'GetUserByID'.
func (s *userService) GetUserByID(ctx context.Context, userID uint) (*domain.User, error) {
	return s.userRepository.GetUserByID(ctx, userID)
}

// Implementation of 'GetUserByMail'.
func (s *userService) GetUserByMail(ctx context.Context, userMail string) (*domain.User, error) {
	return s.userRepository.GetUserByMail(ctx, userMail)
}

// Implementation of 'GetUserByToken'.
func (s *userService) GetUserByToken(ctx context.Context, token string) (*domain.User, error) {
	return s.userRepository.GetUserByToken(ctx, token)
}

// Implementation of 'CreateUser'.
func (s *userService) CreateUser(ctx context.Context, datas *dto.UserInputDTO) (*domain.User, error) {
	return s.userRepository.CreateUser(ctx, datas)
}

// Implementation of 'UpdateUser'.
func (s *userService) UpdateUser(ctx context.Context, user *domain.User, datas *dto.UserInputDTO) error {
	return s.userRepository.UpdateUser(ctx, user, datas)
}

// Implementation of 'DeleteUser'.
func (s *userService) DeleteUser(ctx context.Context, user *domain.User) error {
	return s.userRepository.DeleteUser(ctx, user)
}

// Implementation of 'ResetUser'.
func (s *userService) ResetUser(ctx context.Context, user *domain.User) error {
	return s.userRepository.ResetUser(ctx, user)
}

// Implementation of 'PasswordUser'.
func (s *userService) PasswordUser(ctx context.Context, user *domain.User, pass *dto.PasswordInputDTO) error {
	return s.userRepository.PasswordUser(ctx, user, pass)
}
