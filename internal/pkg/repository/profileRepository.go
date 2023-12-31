package repository

import (
	"context"

	"github.com/raulaguila/go-template/internal/pkg/domain"
	"github.com/raulaguila/go-template/internal/pkg/dto"
	"github.com/raulaguila/go-template/pkg/filter"
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

func (s *profileRepository) applyFilter(ctx context.Context, filter *filter.Filter) *gorm.DB {
	postgres := s.postgres.WithContext(ctx)
	postgres = filter.ApplySearchLike(postgres, "name")
	postgres = filter.ApplyOrder(postgres)

	return postgres
}

func (s *profileRepository) countProfiles(db *gorm.DB) (int64, error) {
	var count int64
	return count, db.Model(&domain.Profile{}).Count(&count).Error
}

func (s *profileRepository) listProfiles(db *gorm.DB) (*[]domain.Profile, error) {
	profiles := &[]domain.Profile{}
	return profiles, db.Preload(clause.Associations).Find(profiles).Error
}

func (s *profileRepository) GetProfilesOutputDTO(ctx context.Context, filter *filter.Filter) (*dto.ItemsOutputDTO, error) {
	postgres := s.applyFilter(ctx, filter)
	count, err := s.countProfiles(postgres)
	if err != nil {
		return nil, err
	}

	postgres = filter.ApplyPagination(postgres)
	items, err := s.listProfiles(postgres)
	if err != nil {
		return nil, err
	}

	return &dto.ItemsOutputDTO{
		Items: items,
		Count: count,
	}, nil
}

func (s *profileRepository) GetProfileByID(ctx context.Context, profileID uint) (*domain.Profile, error) {
	profile := &domain.Profile{}
	return profile, s.postgres.WithContext(ctx).Preload(clause.Associations).First(profile, profileID).Error
}

func (s *profileRepository) CreateProfile(ctx context.Context, datas *dto.ProfileInputDTO) (*domain.Profile, error) {
	profile := &domain.Profile{}
	if err := profile.Bind(datas); err != nil {
		return nil, err
	}

	return profile, s.postgres.WithContext(ctx).Create(profile).Error
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
