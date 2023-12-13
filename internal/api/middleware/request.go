package middleware

import (
	"errors"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/raulaguila/go-template/internal/pkg/domain"
	"github.com/raulaguila/go-template/internal/pkg/i18n"
	"github.com/raulaguila/go-template/internal/pkg/postgre"
	httphelper "github.com/raulaguila/go-template/pkg/http-helper"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func NewObjectMiddleware(postgres *gorm.DB) *ObjectMiddleware {
	return &ObjectMiddleware{
		postgres: postgres,
	}
}

type ObjectMiddleware struct {
	postgres *gorm.DB
}

var ErrInvalidID error = errors.New("invalid id")

const (
	strUser    string = "user"
	strProfile string = "prof"
	strProduct string = "prod"
)

func (ObjectMiddleware) handlerError(c *fiber.Ctx, err error) error {
	translation := c.Locals(httphelper.LocalLang).(*i18n.Translation)

	switch err {
	case ErrInvalidID:
		return httphelper.NewHTTPErrorResponse(c, fiber.StatusBadRequest, translation.ErrInvalidId)
	default:
		log.Println(err.Error())
		return httphelper.NewHTTPErrorResponse(c, fiber.StatusInternalServerError, translation.ErrGeneric)
	}
}

func (s ObjectMiddleware) handlerNotFoundError(c *fiber.Ctx, err error, item string) error {
	translation := c.Locals(httphelper.LocalLang).(*i18n.Translation)

	switch err {
	case gorm.ErrRecordNotFound:
		switch item {
		case strUser:
			return httphelper.NewHTTPErrorResponse(c, fiber.StatusNotFound, translation.ErrUserNotFound)
		case strProfile:
			return httphelper.NewHTTPErrorResponse(c, fiber.StatusNotFound, translation.ErrProfileNotFound)
		case strProduct:
			return httphelper.NewHTTPErrorResponse(c, fiber.StatusNotFound, translation.ErrProductNotFound)
		}
	}

	log.Println(err.Error())
	return s.handlerError(c, err)
}

func (ObjectMiddleware) getID(c *fiber.Ctx) (uint, bool) {
	targetedID, err := c.ParamsInt(httphelper.ParamID, 0)
	return uint(targetedID), (err != nil || targetedID <= 0)
}

func (s *ObjectMiddleware) UserByID(c *fiber.Ctx) error {
	id, err := s.getID(c)
	if err {
		return s.handlerError(c, ErrInvalidID)
	}

	user := &domain.User{}
	if err := s.postgres.WithContext(c.Context()).Preload(postgre.ProfilePermission).First(user, id).Error; err != nil {
		return s.handlerNotFoundError(c, err, strUser)
	}

	c.Locals(httphelper.LocalObject, user)
	return c.Next()
}

func (s *ObjectMiddleware) ProfileByID(c *fiber.Ctx) error {
	id, err := s.getID(c)
	if err {
		return s.handlerError(c, ErrInvalidID)
	}

	profile := &domain.Profile{}
	if err := s.postgres.WithContext(c.Context()).Preload(clause.Associations).First(profile, id).Error; err != nil {
		return s.handlerNotFoundError(c, err, strProfile)
	}

	c.Locals(httphelper.LocalObject, profile)
	return c.Next()
}

func (s *ObjectMiddleware) ProductByID(c *fiber.Ctx) error {
	id, err := s.getID(c)
	if err {
		return s.handlerError(c, ErrInvalidID)
	}

	product := &domain.Product{}
	if err := s.postgres.WithContext(c.Context()).First(product, id).Error; err != nil {
		return s.handlerNotFoundError(c, err, strProduct)
	}

	c.Locals(httphelper.LocalObject, product)
	return c.Next()
}
