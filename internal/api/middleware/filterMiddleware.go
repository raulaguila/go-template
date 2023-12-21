package middleware

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/raulaguila/go-template/internal/pkg/i18n"
	"github.com/raulaguila/go-template/pkg/filter"
	httphelper "github.com/raulaguila/go-template/pkg/http-helper"
)

func getQuery(c *fiber.Ctx, data interface{}) error {
	if err := c.QueryParser(data); err != nil {
		log.Println(err.Error())
		messages := c.Locals(httphelper.LocalLang).(*i18n.Translation)
		return httphelper.NewHTTPErrorResponse(c, fiber.StatusBadRequest, messages.ErrInvalidDatas)
	}

	c.Locals(httphelper.LocalFilter, data)
	return c.Next()
}

func GetGenericFilter(c *fiber.Ctx) error {
	return getQuery(c, filter.NewFilter())
}

func GetUserFilter(c *fiber.Ctx) error {
	return getQuery(c, &filter.UserFilter{
		Filter:    *filter.NewFilter(),
		ProfileID: 0,
	})
}
