package handler

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/raulaguila/go-template/internal/api/middleware"
	"github.com/raulaguila/go-template/internal/pkg/domain"
	"github.com/raulaguila/go-template/internal/pkg/dto"
	"github.com/raulaguila/go-template/internal/pkg/i18n"
	httphelper "github.com/raulaguila/go-template/pkg/http-helper"
	"gorm.io/gorm"
)

type AuthHandler struct {
	authService domain.AuthService
}

func (AuthHandler) handlerError(c *fiber.Ctx, err error) error {
	messages := c.Locals(httphelper.LocalLang).(*i18n.Translation)

	log.Println(err.Error())
	return httphelper.NewHTTPResponse(c, fiber.StatusInternalServerError, messages.ErrGeneric)
}

func (s *AuthHandler) checkCredentials(c *fiber.Ctx) error {
	translation := c.Locals(httphelper.LocalLang).(*i18n.Translation)
	credentials := &dto.LoginInputDTO{}
	if err := c.BodyParser(credentials); err != nil {
		return httphelper.NewHTTPResponse(c, fiber.StatusBadRequest, translation.ErrInvalidDatas)
	}

	user, err := s.authService.GetUserByMail(c.Context(), credentials.Email)
	if err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			return httphelper.NewHTTPResponse(c, fiber.StatusUnauthorized, translation.ErrUserNotFound)
		default:
			log.Println(err.Error())
			return httphelper.NewHTTPResponse(c, fiber.StatusInternalServerError, translation.ErrGeneric)
		}
	}

	if !user.ValidatePassword(credentials.Password) {
		return httphelper.NewHTTPResponse(c, fiber.StatusUnauthorized, translation.ErrIncorrectPassword)
	}

	if !user.Status || user.New {
		return httphelper.NewHTTPResponse(c, fiber.StatusUnauthorized, translation.ErrDisabledUser)
	}

	c.Locals(httphelper.LocalObject, user)
	return c.Next()
}

// Creates a new handler.
func NewAuthHandler(route fiber.Router, as domain.AuthService) {
	handler := &AuthHandler{
		authService: as,
	}

	route.Post("", handler.checkCredentials, handler.login)
	route.Get("", middleware.MidAccess, handler.me)
	route.Put("", middleware.MidRefresh, handler.refresh)
}

// login godoc
// @Summary      User authentication
// @Description  User authentication
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        lang query string false "Language responses"
// @Param        credentials body dto.LoginInputDTO true "Credentials model"
// @Success      200  {object}  domain.AuthResponse
// @Failure      401  {object}  httphelper.HTTPResponse
// @Failure      500  {object}  httphelper.HTTPResponse
// @Router       /auth [post]
func (s *AuthHandler) login(c *fiber.Ctx) error {
	authResponse, err := s.authService.Login(c.Context(), c.Locals(httphelper.LocalObject).(*domain.User))
	if err != nil {
		return s.handlerError(c, err)
	}

	if err := authResponse.Validate(); err != nil {
		return s.handlerError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(authResponse)
}

// me godoc
// @Summary      User authenticated
// @Description  User authenticated
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        Authorization header string false "User token"
// @Param        lang query string false "Language responses"
// @Success      200  {object}  domain.User
// @Failure      401  {object}  httphelper.HTTPResponse
// @Failure      500  {object}  httphelper.HTTPResponse
// @Router       /auth [get]
// @Security	 Bearer
func (s *AuthHandler) me(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(c.Locals(httphelper.LocalUser).(*domain.User))
}

// refresh godoc
// @Summary      User refresh
// @Description  User refresh
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        Authorization header string false "User token"
// @Param        lang query string false "Language responses"
// @Success      200  {object}  domain.TokensResponse
// @Failure      401  {object}  httphelper.HTTPResponse
// @Failure      500  {object}  httphelper.HTTPResponse
// @Router       /auth [put]
func (s *AuthHandler) refresh(c *fiber.Ctx) error {
	tokensResponse, err := s.authService.Refresh(c.Context(), c.Locals(httphelper.LocalUser).(*domain.User))
	if err != nil {
		return s.handlerError(c, err)
	}

	if err := tokensResponse.Validate(); err != nil {
		return s.handlerError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(tokensResponse)
}
