package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/raulaguila/go-template/internal/pkg/domain"
	"github.com/raulaguila/go-template/internal/pkg/dto"
	"github.com/raulaguila/go-template/internal/pkg/postgre"
	gormhelper "github.com/raulaguila/go-template/pkg/gorm-helper"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func NewUserRepository(postgres *gorm.DB) domain.UserRepository {
	return &userRepository{
		postgres: postgres,
	}
}

type userRepository struct {
	postgres *gorm.DB
}

func (s *userRepository) applyFilter(ctx context.Context, filter *gormhelper.UserFilter, pag bool) *gorm.DB {
	postgres := s.postgres.WithContext(ctx)
	if filter.ProfileID != 0 {
		postgres = postgres.Where("profile_id = ?", filter.ProfileID)
	}
	postgres = postgres.Joins(fmt.Sprintf("JOIN %v ON %v.id = %v.profile_id", domain.ProfileTableName, domain.ProfileTableName, domain.UserTableName))
	postgres = filter.ApplySearchLike(postgres, domain.UserTableName+".name", domain.UserTableName+".mail", domain.ProfileTableName+".name")
	postgres = filter.ApplyOrder(postgres)
	if pag {
		postgres = filter.ApplyPagination(postgres)
	}

	return postgres
}

func (s *userRepository) GetUserByID(ctx context.Context, userID uint) (*domain.User, error) {
	user := &domain.User{}
	return user, s.postgres.WithContext(ctx).Preload(postgre.ProfilePermission).First(user, userID).Error
}

func (s *userRepository) GetUsers(ctx context.Context, filter *gormhelper.UserFilter) (*[]domain.User, error) {
	users := &[]domain.User{}
	return users, s.applyFilter(ctx, filter, true).Preload(postgre.ProfilePermission).Find(users).Error
}

func (s *userRepository) CountUsers(ctx context.Context, filter *gormhelper.UserFilter) (int64, error) {
	var count int64
	return count, s.applyFilter(ctx, filter, false).Model(&domain.User{}).Count(&count).Error
}

func (s *userRepository) GetUserByMail(ctx context.Context, mail string) (*domain.User, error) {
	user := &domain.User{Email: mail}
	return user, s.postgres.WithContext(ctx).Preload(postgre.ProfilePermission).Where(user).First(user).Error
}

func (s *userRepository) GetUserByToken(ctx context.Context, token string) (*domain.User, error) {
	user := &domain.User{Token: &token}
	return user, s.postgres.WithContext(ctx).Preload(postgre.ProfilePermission).Where(user).First(user).Error
}

func (s *userRepository) CreateUser(ctx context.Context, datas *dto.UserInputDTO) (uint, error) {
	user := &domain.User{New: true}
	if err := user.Bind(datas); err != nil {
		return 0, err
	}

	if err := s.postgres.WithContext(ctx).Create(user).Error; err != nil {
		return 0, err
	}

	return user.Id, nil
}

func (s *userRepository) UpdateUser(ctx context.Context, user *domain.User, datas *dto.UserInputDTO) error {
	if err := user.Bind(datas); err != nil {
		return err
	}

	return s.postgres.WithContext(ctx).Model(user).Updates(user.ToMap()).Error
}

func (s *userRepository) DeleteUser(ctx context.Context, user *domain.User) error {
	return s.postgres.WithContext(ctx).Delete(user).Error
}

func (s *userRepository) ResetUser(ctx context.Context, user *domain.User) error {
	user.Password = nil
	user.Token = nil
	user.New = true

	return s.UpdateUser(ctx, user, &dto.UserInputDTO{})
}

func (s *userRepository) PasswordUser(ctx context.Context, user *domain.User, pass *dto.PasswordInputDTO) error {
	user.New = false
	user.Token = new(string)
	*user.Token = uuid.New().String()

	hash, err := bcrypt.GenerateFromPassword([]byte(*pass.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = new(string)
	*user.Password = string(hash)

	return s.UpdateUser(ctx, user, &dto.UserInputDTO{})
}
