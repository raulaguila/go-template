package service

import (
	"context"

	"github.com/raulaguila/go-template/internal/pkg/domain"
	"github.com/raulaguila/go-template/internal/pkg/dto"
	"github.com/raulaguila/go-template/pkg/filter"
)

func NewProfileService(r domain.ProfileRepository) domain.ProfileService {
	return &profileService{
		profileRepository: r,
	}
}

type profileService struct {
	profileRepository domain.ProfileRepository
}

// Implementation of 'GetProfileByID'.
func (s *profileService) GetProfileByID(ctx context.Context, profileID uint) (*domain.Profile, error) {
	return s.profileRepository.GetProfileByID(ctx, profileID)
}

// Implementation of 'GetProfiles'.
func (s *profileService) GetProfiles(ctx context.Context, filter *filter.Filter) (*[]domain.Profile, error) {
	return s.profileRepository.GetProfiles(ctx, filter)
}

// Implementation of 'CountProfiles'.
func (s *profileService) CountProfiles(ctx context.Context, filter *filter.Filter) (int64, error) {
	return s.profileRepository.CountProfiles(ctx, filter)
}

// Implementation of 'CreateProfile'.
func (s *profileService) CreateProfile(ctx context.Context, datas *dto.ProfileInputDTO) (uint, error) {
	return s.profileRepository.CreateProfile(ctx, datas)
}

// Implementation of 'UpdateProfile'.
func (s *profileService) UpdateProfile(ctx context.Context, profile *domain.Profile, datas *dto.ProfileInputDTO) error {
	return s.profileRepository.UpdateProfile(ctx, profile, datas)
}

// Implementation of 'DeleteProfile'.
func (s *profileService) DeleteProfile(ctx context.Context, profile *domain.Profile) error {
	return s.profileRepository.DeleteProfile(ctx, profile)
}
