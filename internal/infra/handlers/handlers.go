package handlers

import (
	"log"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"

	"github.com/raulaguila/go-template/docs"
	"github.com/raulaguila/go-template/internal/api/handler"
	"github.com/raulaguila/go-template/internal/api/middleware"
	"github.com/raulaguila/go-template/internal/api/service"
	"github.com/raulaguila/go-template/internal/pkg/i18n"
	"github.com/raulaguila/go-template/internal/pkg/repository"
	httphelper "github.com/raulaguila/go-template/pkg/http-helper"

	"gorm.io/gorm"
)

func HandleRequests(app *fiber.App, postgresdb *gorm.DB) {
	if strings.ToLower(os.Getenv("API_SWAGGO")) == "true" {
		docs.SwaggerInfo.Version = os.Getenv("SYS_VERSION")

		// 	// Config swagger
		app.Get("/swagger/*", swagger.New(swagger.Config{
			DisplayRequestDuration: true,
			DocExpansion:           "none",
			ValidatorUrl:           "none",
		}))
	}

	// Create repositories.
	profileRepository := repository.NewProfileRepository(postgresdb)
	userRepository := repository.NewUserRepository(postgresdb)
	authRepository := repository.NewAuthRepository(userRepository)
	productRepository := repository.NewProductRepository(postgresdb)

	// Create services.
	profileService := service.NewProfileService(profileRepository)
	userService := service.NewUserService(userRepository)
	authService := service.NewAuthService(authRepository)
	productService := service.NewProductService(productRepository)

	objMiddleware := middleware.NewObjectMiddleware(postgresdb)

	middleware.MidAccess = middleware.Auth(os.Getenv("ACCESS_TOKEN_PUBLIC"), authService)
	middleware.MidRefresh = middleware.Auth(os.Getenv("RFRESH_TOKEN_PUBLIC"), authService)

	// Prepare endpoints for the API.
	handler.NewMiscHandler(app.Group(""))
	handler.NewProfileHandler(app.Group("/profile"), profileService, objMiddleware)
	handler.NewUserHandler(app.Group("/user"), userService, objMiddleware)
	handler.NewAuthHandler(app.Group("/auth"), authService)
	handler.NewProductHandler(app.Group("/product"), productService, objMiddleware)

	// Prepare an endpoint for 'Not Found'.
	app.All("*", func(c *fiber.Ctx) error {
		messages := c.Locals(httphelper.LocalLang).(*i18n.Translation)
		return httphelper.NewHTTPErrorResponse(c, fiber.StatusNotFound, messages.ErrorNonexistentRoute)
	})

	log.Fatal(app.Listen(":" + os.Getenv("API_PORT")))
}
