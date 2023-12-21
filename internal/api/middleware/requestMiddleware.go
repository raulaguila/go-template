package middleware

import (
	"errors"
	"log"
	"reflect"

	"github.com/gofiber/fiber/v2"
	"github.com/raulaguila/go-template/internal/pkg/domain"
	"github.com/raulaguila/go-template/internal/pkg/i18n"
	httphelper "github.com/raulaguila/go-template/pkg/http-helper"
	"gorm.io/gorm"
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

func (RequesttMiddleware) handlerError(c *fiber.Ctx, err error) error {
	translation := c.Locals(httphelper.LocalLang).(*i18n.Translation)

	switch err {
	case ErrInvalidID:
		return httphelper.NewHTTPErrorResponse(c, fiber.StatusBadRequest, translation.ErrInvalidId)
	default:
		log.Println(err.Error())
		return httphelper.NewHTTPErrorResponse(c, fiber.StatusInternalServerError, translation.ErrGeneric)
	}
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

	return s.handlerError(c, err)
}

func (RequesttMiddleware) getID(c *fiber.Ctx) (uint, bool) {
	targetedID, err := c.ParamsInt(httphelper.ParamID, 0)
	return uint(targetedID), (err != nil || targetedID <= 0)
}

func (s *RequesttMiddleware) ItemByID(item interface{}, itemType string, preload ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		valuesData := reflect.ValueOf(item).Elem()
		valuesData.Set(reflect.Zero(valuesData.Type()))

		id, err := s.getID(c)
		if err {
			return s.handlerError(c, ErrInvalidID)
		}

		db := s.postgres.WithContext(c.Context())
		for _, item := range preload {
			db = db.Preload(item)
		}

		if err := db.First(item, id).Error; err != nil {
			return s.handlerDBError(c, err, itemType)
		}

		c.Locals(httphelper.LocalObject, item)
		return c.Next()
	}
}
