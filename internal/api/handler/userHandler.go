package handler

import (
	"errors"
	"log"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/raulaguila/go-template/internal/api/middleware"
	"github.com/raulaguila/go-template/internal/pkg/domain"
	"github.com/raulaguila/go-template/internal/pkg/dto"
	"github.com/raulaguila/go-template/internal/pkg/i18n"
	"github.com/raulaguila/go-template/internal/pkg/postgre"
	"github.com/raulaguila/go-template/pkg/filter"
	httphelper "github.com/raulaguila/go-template/pkg/http-helper"
	"github.com/raulaguila/go-template/pkg/pgerror"
	"github.com/raulaguila/go-template/pkg/validator"
	"gorm.io/gorm"
)

type UserHandler struct {
	userService domain.UserService
}

func (s UserHandler) foreignKeyViolatedFrom(c *fiber.Ctx, messages *i18n.Translation) error {
	switch c.Method() {
	case fiber.MethodPut, fiber.MethodPost, fiber.MethodPatch:
		return httphelper.NewHTTPErrorResponse(c, fiber.StatusBadRequest, messages.ErrProfileNotFound)
	case fiber.MethodDelete:
		return httphelper.NewHTTPErrorResponse(c, fiber.StatusBadRequest, messages.ErrUserUsed)
	default:
		return httphelper.NewHTTPErrorResponse(c, fiber.StatusInternalServerError, messages.ErrGeneric)
	}
}

func (s UserHandler) handlerError(c *fiber.Ctx, err error) error {
	messages := c.Locals(httphelper.LocalLang).(*i18n.Translation)

	switch pgerror.HandlerError(err) {
	case pgerror.ErrDuplicatedKey:
		return httphelper.NewHTTPErrorResponse(c, fiber.StatusConflict, messages.ErrUserRegistered)
	case pgerror.ErrForeignKeyViolated:
		return s.foreignKeyViolatedFrom(c, messages)
	case pgerror.ErrUndefinedColumn:
		return httphelper.NewHTTPErrorResponse(c, fiber.StatusBadRequest, messages.ErrUndefinedColumn)
	}

	if errors.As(err, &validator.ErrValidator) {
		return httphelper.NewHTTPErrorResponse(c, fiber.StatusBadRequest, err)
	}

	log.Println(err.Error())
	return httphelper.NewHTTPErrorResponse(c, fiber.StatusInternalServerError, messages.ErrGeneric)
}

func (h *UserHandler) existUserByEmail(c *fiber.Ctx) error {
	translation := c.Locals(httphelper.LocalLang).(*i18n.Translation)

	mail := strings.ReplaceAll(c.Params(httphelper.ParamMail), "%40", "@")
	user, err := h.userService.GetUserByMail(c.Context(), mail)
	if err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			return httphelper.NewHTTPErrorResponse(c, fiber.StatusNotFound, translation.ErrUserNotFound)
		default:
			log.Println(err.Error())
			return httphelper.NewHTTPErrorResponse(c, fiber.StatusInternalServerError, translation.ErrGeneric)
		}
	}

	c.Locals(httphelper.LocalObject, user)
	return c.Next()
}

// Creates a new handler.
func NewUserHandler(route fiber.Router, us domain.UserService, mid *middleware.RequesttMiddleware) {
	handler := &UserHandler{
		userService: us,
	}

	route.Patch("/:"+httphelper.ParamMail+"/passw", handler.existUserByEmail, middleware.GetDTO(&dto.PasswordInputDTO{}), handler.passwordUser)

	route.Use(middleware.MidAccess)
	userByID := mid.ItemByID(&domain.User{}, domain.UserTableName, postgre.ProfilePermission)

	route.Get("", middleware.GetUserFilter, handler.getUsers)
	route.Post("", middleware.GetDTO(&dto.UserInputDTO{}), handler.createUser)
	route.Get("/:"+httphelper.ParamID, userByID, handler.getUser)
	route.Put("/:"+httphelper.ParamID, userByID, middleware.GetDTO(&dto.UserInputDTO{}), handler.updateUser)
	route.Delete("/:"+httphelper.ParamID, userByID, handler.deleteUser)
	route.Patch("/:"+httphelper.ParamID+"/reset", userByID, handler.resetUser)
}

// getUsers godoc
// @Summary      Get users
// @Description  Get all users
// @Tags         User
// @Accept       json
// @Produce      json
// @Param        lang query string false "Language responses"
// @Param        filter query filter.UserFilter false "Optional Filter"
// @Success      200  {array}   dto.ItemsOutputDTO
// @Failure      500  {object}  httphelper.HTTPResponse
// @Router       /user [get]
// @Security	 Bearer
func (h *UserHandler) getUsers(c *fiber.Ctx) error {
	response, err := h.userService.GetUsersOutputDTO(c.Context(), c.Locals(httphelper.LocalFilter).(*filter.UserFilter))
	if err != nil {
		return h.handlerError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(response)
}

// createUser godoc
// @Summary      Insert user
// @Description  Insert user
// @Tags         User
// @Accept       json
// @Produce      json
// @Param        lang query string false "Language responses"
// @Param        user body dto.UserInputDTO true "User model"
// @Success      201  {object}  domain.User
// @Failure      400  {object}  httphelper.HTTPResponse
// @Failure      409  {object}  httphelper.HTTPResponse
// @Failure      500  {object}  httphelper.HTTPResponse
// @Router       /user [post]
// @Security	 Bearer
func (h *UserHandler) createUser(c *fiber.Ctx) error {
	userDTO := c.Locals(httphelper.LocalDTO).(*dto.UserInputDTO)
	user, err := h.userService.CreateUser(c.Context(), userDTO)
	if err != nil {
		return h.handlerError(c, err)
	}

	return c.Status(fiber.StatusCreated).JSON(user)
}

// getUser godoc
// @Summary      Get user
// @Description  Get user by ID
// @Tags         User
// @Accept       json
// @Produce      json
// @Param        lang query string false "Language responses"
// @Param        id     path    int     true        "User ID"
// @Success      200  {object}  domain.User
// @Failure      400  {object}  httphelper.HTTPResponse
// @Failure      404  {object}  httphelper.HTTPResponse
// @Failure      500  {object}  httphelper.HTTPResponse
// @Router       /user/{id} [get]
// @Security	 Bearer
func (h *UserHandler) getUser(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(c.Locals(httphelper.LocalObject).(*domain.User))
}

// updateUser godoc
// @Summary      Update user
// @Description  Update user by ID
// @Tags         User
// @Accept       json
// @Produce      json
// @Param        lang query string false "Language responses"
// @Param        id     path    int     true        "User ID"
// @Param        user body dto.UserInputDTO true "User model"
// @Success      200  {object}  domain.User
// @Failure      400  {object}  httphelper.HTTPResponse
// @Failure      404  {object}  httphelper.HTTPResponse
// @Failure      500  {object}  httphelper.HTTPResponse
// @Router       /user/{id} [put]
// @Security	 Bearer
func (h *UserHandler) updateUser(c *fiber.Ctx) error {
	userDTO := c.Locals(httphelper.LocalDTO).(*dto.UserInputDTO)
	user := c.Locals(httphelper.LocalObject).(*domain.User)
	if err := h.userService.UpdateUser(c.Context(), user, userDTO); err != nil {
		return h.handlerError(c, err)
	}

	updated, err := h.userService.GetUserByID(c.Context(), user.Id)
	if err != nil {
		return h.handlerError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(updated)
}

// deleteUser godoc
// @Summary      Delete user
// @Description  Delete user by ID
// @Tags         User
// @Accept       json
// @Produce      json
// @Param        lang query string false "Language responses"
// @Param        id     path    int     true        "User ID"
// @Success      204  {object}  nil
// @Failure      404  {object}  httphelper.HTTPResponse
// @Failure      500  {object}  httphelper.HTTPResponse
// @Router       /user/{id} [delete]
// @Security	 Bearer
func (h *UserHandler) deleteUser(c *fiber.Ctx) error {
	if err := h.userService.DeleteUser(c.Context(), c.Locals(httphelper.LocalObject).(*domain.User)); err != nil {
		return h.handlerError(c, err)
	}

	return c.Status(fiber.StatusNoContent).Send(nil)
}

// resetUser godoc
// @Summary      Reset user password
// @Description  Reset user password by ID
// @Tags         User
// @Accept       json
// @Produce      json
// @Param        lang query string false "Language responses"
// @Param        id     path    int     true        "User ID"
// @Success      200  {object}  domain.User
// @Failure      404  {object}  httphelper.HTTPResponse
// @Failure      500  {object}  httphelper.HTTPResponse
// @Router       /user/{id}/reset [patch]
// @Security	 Bearer
func (h *UserHandler) resetUser(c *fiber.Ctx) error {
	user := c.Locals(httphelper.LocalObject).(*domain.User)

	if !user.New {
		if err := h.userService.ResetUser(c.Context(), user); err != nil {
			return h.handlerError(c, err)
		}
		updated, err := h.userService.GetUserByID(c.Context(), user.Id)
		if err != nil {
			return h.handlerError(c, err)
		}
		return c.Status(fiber.StatusOK).JSON(updated)
	}

	return c.Status(fiber.StatusOK).JSON(user)
}

// passwordUser godoc
// @Summary      Set user password
// @Description  Set user password by ID
// @Tags         User
// @Accept       json
// @Produce      json
// @Param        lang query string false "Language responses"
// @Param        email     path    string     true        "User email"
// @Param        password body dto.PasswordInputDTO true "Password model"
// @Success      200  {object}  domain.User
// @Failure      404  {object}  httphelper.HTTPResponse
// @Failure      500  {object}  httphelper.HTTPResponse
// @Router       /user/{email}/passw [patch]
func (h *UserHandler) passwordUser(c *fiber.Ctx) error {
	pass := c.Locals(httphelper.LocalDTO).(*dto.PasswordInputDTO)
	if !pass.IsValid() {
		messages := c.Locals(httphelper.LocalLang).(*i18n.Translation)
		return httphelper.NewHTTPErrorResponse(c, fiber.StatusBadRequest, messages.ErrPassUnmatch)
	}

	user := c.Locals(httphelper.LocalObject).(*domain.User)
	if !user.New {
		messages := c.Locals(httphelper.LocalLang).(*i18n.Translation)
		return httphelper.NewHTTPErrorResponse(c, fiber.StatusBadRequest, messages.ErrUserHasPass)
	}

	if err := h.userService.PasswordUser(c.Context(), user, pass); err != nil {
		return h.handlerError(c, err)
	}

	updated, err := h.userService.GetUserByID(c.Context(), user.Id)
	if err != nil {
		return h.handlerError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(updated)
}
