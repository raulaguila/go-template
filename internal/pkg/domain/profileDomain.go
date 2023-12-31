package domain

import (
	"context"

	"github.com/raulaguila/go-template/internal/pkg/dto"
	"github.com/raulaguila/go-template/pkg/filter"
	"github.com/raulaguila/go-template/pkg/validator"
)

const ProfileTableName string = "profiles"

type (
	Permissions struct {
		Id            uint `json:"-" gorm:"primarykey"`
		ProfileID     uint `json:"-" gorm:"column:profile_id;unique;"`
		UserModule    bool `json:"user_module" gorm:"column:user;type:bool;not null;"`
		ProfileModule bool `json:"profile_module" gorm:"column:profile;type:bool;not null;"`
		ProductModule bool `json:"product_module" gorm:"column:product;type:bool;not null;"`
	}

	Profile struct {
		Base
		Name        string      `json:"name" gorm:"column:name;type:varchar(100);unique;not null;" validate:"required,min=4"`
		Permissions Permissions `json:"permissions" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	}

	ProfileRepository interface {
		GetProfileByID(context.Context, uint) (*Profile, error)
		GetProfilesOutputDTO(context.Context, *filter.Filter) (*dto.ItemsOutputDTO, error)
		CreateProfile(context.Context, *dto.ProfileInputDTO) (*Profile, error)
		UpdateProfile(context.Context, *Profile, *dto.ProfileInputDTO) error
		DeleteProfile(context.Context, *Profile) error
	}

	ProfileService interface {
		GetProfileByID(context.Context, uint) (*Profile, error)
		GetProfilesOutputDTO(context.Context, *filter.Filter) (*dto.ItemsOutputDTO, error)
		CreateProfile(context.Context, *dto.ProfileInputDTO) (*Profile, error)
		UpdateProfile(context.Context, *Profile, *dto.ProfileInputDTO) error
		DeleteProfile(context.Context, *Profile) error
	}
)

func (Profile) TableName() string {
	return ProfileTableName
}

func (s *Permissions) ToMap() *map[string]interface{} {
	return &map[string]interface{}{
		"user":    s.UserModule,
		"profile": s.ProfileModule,
		"product": s.ProductModule,
	}
}

func (s *Profile) Bind(p *dto.ProfileInputDTO) error {
	if p.Name != nil {
		s.Name = *p.Name
	}

	if p.Permissions.UserModule != nil {
		s.Permissions.UserModule = *p.Permissions.UserModule
	}

	if p.Permissions.ProfileModule != nil {
		s.Permissions.ProfileModule = *p.Permissions.ProfileModule
	}

	if p.Permissions.ProductModule != nil {
		s.Permissions.ProductModule = *p.Permissions.ProductModule
	}

	return validator.StructValidator.Validate(s)
}
