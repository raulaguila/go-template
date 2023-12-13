package handler

import (
	"errors"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/raulaguila/go-template/internal/api/middleware"
	"github.com/raulaguila/go-template/internal/pkg/domain"
	"github.com/raulaguila/go-template/internal/pkg/dto"
	"github.com/raulaguila/go-template/internal/pkg/i18n"
	gormhelper "github.com/raulaguila/go-template/pkg/gorm-helper"
	httphelper "github.com/raulaguila/go-template/pkg/http-helper"
	"github.com/raulaguila/go-template/pkg/postgresql"
	"github.com/raulaguila/go-template/pkg/validator"
)

type ProfileHandler struct {
	profileService domain.ProfileService
}

func (ProfileHandler) handlerError(c *fiber.Ctx, err error) error {
	messages := c.Locals(httphelper.LocalLang).(*i18n.Translation)

	switch postgresql.HandlerError(err) {
	case postgresql.ErrDuplicatedKey:
		return httphelper.NewHTTPErrorResponse(c, fiber.StatusConflict, messages.ErrProfileRegistered)
	case postgresql.ErrForeignKeyViolated:
		return httphelper.NewHTTPErrorResponse(c, fiber.StatusBadRequest, messages.ErrProfileUsed)
	case postgresql.ErrUndefinedColumn:
		return httphelper.NewHTTPErrorResponse(c, fiber.StatusBadRequest, messages.ErrUndefinedColumn)
	}

	if errors.As(err, &validator.ErrValidator) {
		return httphelper.NewHTTPErrorResponse(c, fiber.StatusBadRequest, err)
	}

	log.Println(err.Error())
	return httphelper.NewHTTPErrorResponse(c, fiber.StatusInternalServerError, messages.ErrGeneric)
}

// Creates a new handler.
func NewProfileHandler(route fiber.Router, ps domain.ProfileService, mid *middleware.ObjectMiddleware) {
	handler := &ProfileHandler{
		profileService: ps,
	}

	route.Use(middleware.MidAccess)

	route.Get("", middleware.GetGenericFilter, handler.getProfiles)
	route.Post("", middleware.GetDTO(&dto.ProfileInputDTO{}), handler.createProfile)
	route.Get("/:"+httphelper.ParamID, mid.ProfileByID, handler.getProfile)
	route.Put("/:"+httphelper.ParamID, mid.ProfileByID, middleware.GetDTO(&dto.ProfileInputDTO{}), handler.updateProfile)
	route.Delete("/:"+httphelper.ParamID, mid.ProfileByID, handler.deleteProfile)
}

// getProfiles godoc
// @Summary      Get profiles
// @Description  Get profiles
// @Tags         Profile
// @Accept       json
// @Produce      json
// @Param        lang query string false "Language responses"
// @Param        filter query gormhelper.Filter false "Optional Filter"
// @Success      200  {array}   dto.ItemsOutputDTO
// @Failure      500  {object}  httphelper.HTTPResponse
// @Router       /profile [get]
// @Security	 Bearer
func (h *ProfileHandler) getProfiles(c *fiber.Ctx) error {
	profiles, err := h.profileService.GetProfiles(c.Context(), c.Locals(httphelper.LocalFilter).(*gormhelper.Filter))
	if err != nil {
		return h.handlerError(c, err)
	}

	count, err := h.profileService.CountProfiles(c.Context(), c.Locals(httphelper.LocalFilter).(*gormhelper.Filter))
	if err != nil {
		return h.handlerError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(&dto.ItemsOutputDTO{
		Items: profiles,
		Count: count,
	})
}

// createProfile godoc
// @Summary      Insert profile
// @Description  Insert profile
// @Tags         Profile
// @Accept       json
// @Produce      json
// @Param        lang query string false "Language responses"
// @Param        profile body dto.ProfileInputDTO true "Profile model"
// @Success      201  {object}  domain.Profile
// @Failure      400  {object}  httphelper.HTTPResponse
// @Failure      409  {object}  httphelper.HTTPResponse
// @Failure      500  {object}  httphelper.HTTPResponse
// @Router       /profile [post]
// @Security	 Bearer
func (h *ProfileHandler) createProfile(c *fiber.Ctx) error {
	id, err := h.profileService.CreateProfile(c.Context(), c.Locals(httphelper.LocalDTO).(*dto.ProfileInputDTO))
	if err != nil {
		return h.handlerError(c, err)
	}

	profile, err := h.profileService.GetProfileByID(c.Context(), id)
	if err != nil {
		return h.handlerError(c, err)
	}

	return c.Status(fiber.StatusCreated).JSON(profile)
}

// getProfile godoc
// @Summary      Get profile by ID
// @Description  Get profile by ID
// @Tags         Profile
// @Accept       json
// @Produce      json
// @Param        lang query string false "Language responses"
// @Param        id     path    int     true        "Profile ID"
// @Success      200  {object}  domain.Profile
// @Failure      400  {object}  httphelper.HTTPResponse
// @Failure      404  {object}  httphelper.HTTPResponse
// @Failure      500  {object}  httphelper.HTTPResponse
// @Router       /profile/{id} [get]
// @Security	 Bearer
func (h *ProfileHandler) getProfile(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(c.Locals(httphelper.LocalObject).(*domain.Profile))
}

// updateProfile godoc
// @Summary      Update profile
// @Description  Update profile by ID
// @Tags         Profile
// @Accept       json
// @Produce      json
// @Param        lang query string false "Language responses"
// @Param        id     path    int     true        "Profile ID"
// @Param        profile body dto.ProfileInputDTO true "Profile model"
// @Success      200  {object}  domain.Profile
// @Failure      400  {object}  httphelper.HTTPResponse
// @Failure      404  {object}  httphelper.HTTPResponse
// @Failure      500  {object}  httphelper.HTTPResponse
// @Router       /profile/{id} [put]
// @Security	 Bearer
func (h *ProfileHandler) updateProfile(c *fiber.Ctx) error {
	profile := c.Locals(httphelper.LocalObject).(*domain.Profile)
	if err := h.profileService.UpdateProfile(c.Context(), profile, c.Locals(httphelper.LocalDTO).(*dto.ProfileInputDTO)); err != nil {
		return h.handlerError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(profile)
}

// deleteProfile godoc
// @Summary      Delete profile
// @Description  Delete profile by ID
// @Tags         Profile
// @Accept       json
// @Produce      json
// @Param        lang query string false "Language responses"
// @Param        id     path    int     true        "Profile ID"
// @Success      204  {object}  nil
// @Failure      404  {object}  httphelper.HTTPResponse
// @Failure      500  {object}  httphelper.HTTPResponse
// @Router       /profile/{id} [delete]
// @Security	 Bearer
func (h *ProfileHandler) deleteProfile(c *fiber.Ctx) error {
	if err := h.profileService.DeleteProfile(c.Context(), c.Locals(httphelper.LocalObject).(*domain.Profile)); err != nil {
		return h.handlerError(c, err)
	}

	return c.Status(fiber.StatusNoContent).Send(nil)
}
