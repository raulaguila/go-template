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

func NewRequesttMiddleware(postgres *gorm.DB) *RequesttMiddleware {
	return &RequesttMiddleware{
		postgres: postgres,
	}
}

type RequesttMiddleware struct {
	postgres *gorm.DB
}

var ErrInvalidID error = errors.New("invalid id")

func (RequesttMiddleware) handlerError(c *fiber.Ctx, err error, translation *i18n.Translation) error {
	switch err {
	case ErrInvalidID:
		return httphelper.NewHTTPErrorResponse(c, fiber.StatusBadRequest, translation.ErrInvalidId)
	}

	log.Println(err.Error())
	return httphelper.NewHTTPErrorResponse(c, fiber.StatusInternalServerError, translation.ErrGeneric)
}

func (s RequesttMiddleware) handlerDBError(c *fiber.Ctx, err error, item string) error {
	translation := c.Locals(httphelper.LocalLang).(*i18n.Translation)

	switch err {
	case gorm.ErrRecordNotFound:
		switch item {
		case domain.UserTableName:
			return httphelper.NewHTTPErrorResponse(c, fiber.StatusNotFound, translation.ErrUserNotFound)
		case domain.ProfileTableName:
			return httphelper.NewHTTPErrorResponse(c, fiber.StatusNotFound, translation.ErrProfileNotFound)
		case domain.ProductTableName:
			return httphelper.NewHTTPErrorResponse(c, fiber.StatusNotFound, translation.ErrProductNotFound)
		}
	}

	return s.handlerError(c, err, translation)
}

func (s *RequesttMiddleware) itemByID(c *fiber.Ctx, item interface{}, itemType string, preload ...string) error {
	id, err := c.ParamsInt(httphelper.ParamID, 0)
	if err != nil || id < 1 {
		translation := c.Locals(httphelper.LocalLang).(*i18n.Translation)
		return s.handlerError(c, ErrInvalidID, translation)
	}

	postgres := s.postgres.WithContext(c.Context())
	for _, pre := range preload {
		postgres = postgres.Preload(pre)
	}

	if err := postgres.First(item, id).Error; err != nil {
		return s.handlerDBError(c, err, itemType)
	}

	c.Locals(httphelper.LocalObject, item)
	return c.Next()
}

func (s *RequesttMiddleware) ProfileByID(c *fiber.Ctx) error {
	return s.itemByID(c, &domain.Profile{}, domain.ProfileTableName, clause.Associations)
}

func (s *RequesttMiddleware) UserByID(c *fiber.Ctx) error {
	return s.itemByID(c, &domain.User{}, domain.UserTableName, postgre.ProfilePermission)
}

func (s *RequesttMiddleware) ProductByID(c *fiber.Ctx) error {
	return s.itemByID(c, &domain.Product{}, domain.ProductTableName)
}
