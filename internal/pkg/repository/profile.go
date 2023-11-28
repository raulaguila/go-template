package repository

import (
	"context"

	"github.com/raulaguila/go-template/internal/pkg/domain"
	"github.com/raulaguila/go-template/internal/pkg/dto"
	gormhelper "github.com/raulaguila/go-template/pkg/gorm-helper"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func NewProfileRepository(postgres *gorm.DB) domain.ProfileRepository {
	return &profileRepository{
		postgres: postgres,
	}
}

type profileRepository struct {
	postgres *gorm.DB
}

func (s *profileRepository) applyFilter(ctx context.Context, filter *gormhelper.Filter, pag bool) *gorm.DB {
	postgres := s.postgres.WithContext(ctx)
	postgres = filter.ApplySearchLike(postgres, "name")
	postgres = filter.ApplyOrder(postgres)
	if pag {
		postgres = filter.ApplyPagination(postgres)
	}

	return postgres
}

func (s *profileRepository) GetProfileByID(ctx context.Context, profileID uint) (*domain.Profile, error) {
	profile := &domain.Profile{}
	return profile, s.postgres.WithContext(ctx).Preload(clause.Associations).First(profile, profileID).Error
}

func (s *profileRepository) GetProfiles(ctx context.Context, filter *gormhelper.Filter) (*[]domain.Profile, error) {
	profiles := &[]domain.Profile{}
	return profiles, s.applyFilter(ctx, filter, true).Preload(clause.Associations).Find(profiles).Error
}

func (s *profileRepository) CountProfiles(ctx context.Context, filter *gormhelper.Filter) (int64, error) {
	var count int64
	return count, s.applyFilter(ctx, filter, false).Model(&domain.Profile{}).Count(&count).Error
}

func (s *profileRepository) CreateProfile(ctx context.Context, datas *dto.ProfileInputDTO) (uint, error) {
	profile := &domain.Profile{}
	if err := profile.Bind(datas); err != nil {
		return 0, err
	}

	return profile.Id, s.postgres.WithContext(ctx).Create(profile).Error
}

func (s *profileRepository) UpdateProfile(ctx context.Context, profile *domain.Profile, datas *dto.ProfileInputDTO) error {
	if err := profile.Bind(datas); err != nil {
		return err
	}

	if err := s.postgres.WithContext(ctx).Model(&profile.Permissions).Updates(profile.Permissions.ToMap()).Error; err != nil {
		return err
	}

	return s.postgres.WithContext(ctx).Updates(profile).Error
}

func (s *profileRepository) DeleteProfile(ctx context.Context, profile *domain.Profile) error {
	return s.postgres.WithContext(ctx).Select(clause.Associations).Delete(profile).Error
}
