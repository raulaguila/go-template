package database

import (
	"context"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/raulaguila/go-template/internal/pkg/domain"
	"github.com/raulaguila/go-template/pkg/helpers"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func autoMigrate(postgresdb *gorm.DB) {
	helpers.PanicIfErr(postgresdb.AutoMigrate(&domain.Permissions{}))
	helpers.PanicIfErr(postgresdb.AutoMigrate(&domain.Profile{}))
	helpers.PanicIfErr(postgresdb.AutoMigrate(&domain.User{}))
	helpers.PanicIfErr(postgresdb.AutoMigrate(&domain.Product{}))
}

func createDefaults(postgresdb *gorm.DB) {
	profile := &domain.Profile{
		Name: "ROOT",
		Permissions: domain.Permissions{
			UserModule:    true,
			ProfileModule: true,
			ProductModule: true,
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	helpers.PanicIfErr(postgresdb.WithContext(ctx).FirstOrCreate(profile, "name = ?", profile.Name).Error)

	user := &domain.User{
		Name:      os.Getenv("ADM_NAME"),
		Email:     os.Getenv("ADM_MAIL"),
		Status:    true,
		ProfileID: profile.Id,
		New:       false,
		Token:     new(string),
		Password:  new(string),
	}

	token := uuid.New().String()
	*user.Token = token

	hash, err := bcrypt.GenerateFromPassword([]byte(os.Getenv("ADM_PASS")), bcrypt.DefaultCost)
	helpers.PanicIfErr(err)
	user.Password = new(string)
	*user.Password = string(hash)

	helpers.PanicIfErr(postgresdb.WithContext(ctx).FirstOrCreate(user, "mail = ?", user.Email).Error)
}
