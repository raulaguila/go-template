package middleware

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/raulaguila/go-template/internal/pkg/i18n"
	gormhelper "github.com/raulaguila/go-template/pkg/gorm-helper"
	httphelper "github.com/raulaguila/go-template/pkg/http-helper"
)

func getFilter(c *fiber.Ctx, data interface{}) error {
	if err := c.QueryParser(data); err != nil {
		log.Println(err.Error())
		messages := c.Locals(httphelper.LocalLang).(*i18n.Translation)
		return httphelper.NewHTTPError(c, fiber.StatusBadRequest, messages.ErrInvalidDatas)
	}

	c.Locals(httphelper.LocalFilter, data)
	return c.Next()
}

func GetGenericFilter(c *fiber.Ctx) error {
	return getFilter(c, gormhelper.NewFilter())
}

func GetUserFilter(c *fiber.Ctx) error {
	return getFilter(c, &gormhelper.UserFilter{
		Filter:    *gormhelper.NewFilter(),
		ProfileID: 0,
	})
}
