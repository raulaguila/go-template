package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/keyauth"
	"github.com/raulaguila/go-template/internal/pkg/domain"
	"github.com/raulaguila/go-template/internal/pkg/i18n"
	httphelper "github.com/raulaguila/go-template/pkg/http-helper"
)

var (
	MidAccess  fiber.Handler
	MidRefresh fiber.Handler
)

func Auth(base64key string, ar domain.AuthService) fiber.Handler {
	return keyauth.New(keyauth.Config{
		KeyLookup:  "header:" + fiber.HeaderAuthorization,
		AuthScheme: "Bearer",
		ContextKey: "token",
		Next: func(c *fiber.Ctx) bool {
			// Filter request to skip middleware
			// true to skip, false to not skip
			return false
		},
		SuccessHandler: func(c *fiber.Ctx) error {
			return c.Next()
		},
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return httphelper.NewHTTPError(c, fiber.StatusUnauthorized, err)
		},
		Validator: func(c *fiber.Ctx, key string) (bool, error) {
			user, err := ar.Me(c.Context(), key, base64key)
			if err != nil || !user.Status {
				translation := c.Locals(httphelper.LocalLang).(*i18n.Translation)
				if err != nil {
					return false, translation.ErrExpiredToken
				}
				return false, translation.ErrDisabledUser
			}
			c.Locals(httphelper.LocalUser, user)
			return true, nil
		},
	})
}
