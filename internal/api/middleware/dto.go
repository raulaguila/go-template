package middleware

import (
	"log"
	"reflect"

	"github.com/gofiber/fiber/v2"
	"github.com/raulaguila/go-template/internal/pkg/i18n"
	httphelper "github.com/raulaguila/go-template/pkg/http-helper"
)

func GetDTO(data interface{}) fiber.Handler {
	return func(c *fiber.Ctx) error {
		valuesData := reflect.ValueOf(data).Elem()
		valuesData.Set(reflect.Zero(valuesData.Type()))

		if err := c.BodyParser(data); err != nil {
			log.Println(err.Error())
			messages := c.Locals(httphelper.LocalLang).(*i18n.Translation)
			return httphelper.NewHTTPError(c, fiber.StatusBadRequest, messages.ErrInvalidDatas)
		}

		c.Locals(httphelper.LocalDTO, data)
		return c.Next()
	}
}
