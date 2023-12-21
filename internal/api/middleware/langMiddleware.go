package middleware

import (
	"os"
	"slices"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/raulaguila/go-template/internal/pkg/i18n"
	httphelper "github.com/raulaguila/go-template/pkg/http-helper"
)

func GetRequestLanguage(c *fiber.Ctx) error {
	lang := strings.ToLower(c.Query("lang", os.Getenv("SYS_LANGUAGE")))[:2]

	if !slices.Contains(strings.Split(os.Getenv("SYS_LANGUAGES"), ","), lang) {
		lang = os.Getenv("SYS_LANGUAGE")
	}

	c.Locals(httphelper.LocalLang, i18n.I18nTranslations[lang])
	return c.Next()
}
