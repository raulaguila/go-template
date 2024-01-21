package handler

import (
	"errors"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/raulaguila/go-template/internal/api/middleware"
	"github.com/raulaguila/go-template/internal/pkg/domain"
	"github.com/raulaguila/go-template/internal/pkg/dto"
	"github.com/raulaguila/go-template/internal/pkg/i18n"
	"github.com/raulaguila/go-template/pkg/filter"
	httphelper "github.com/raulaguila/go-template/pkg/http-helper"
	"github.com/raulaguila/go-template/pkg/pgerror"
	"github.com/raulaguila/go-template/pkg/validator"
)

type ProfileHandler struct {
	profileService domain.ProfileService
}

func (ProfileHandler) foreignKeyViolatedMethod(c *fiber.Ctx, translation *i18n.Translation) error {
	switch c.Method() {
	case fiber.MethodPut, fiber.MethodPost, fiber.MethodPatch:
		return httphelper.NewHTTPResponse(c, fiber.StatusBadRequest, translation.ErrProfileNotFound)
	case fiber.MethodDelete:
		return httphelper.NewHTTPResponse(c, fiber.StatusBadRequest, translation.ErrProfileUsed)
	default:
		return httphelper.NewHTTPResponse(c, fiber.StatusInternalServerError, translation.ErrGeneric)
	}
}

func (s ProfileHandler) handlerError(c *fiber.Ctx, err error) error {
	messages := c.Locals(httphelper.LocalLang).(*i18n.Translation)

	switch pgerror.HandlerError(err) {
	case pgerror.ErrDuplicatedKey:
		return httphelper.NewHTTPResponse(c, fiber.StatusConflict, messages.ErrProfileRegistered)
	case pgerror.ErrForeignKeyViolated:
		return s.foreignKeyViolatedMethod(c, messages)
	case pgerror.ErrUndefinedColumn:
		return httphelper.NewHTTPResponse(c, fiber.StatusBadRequest, messages.ErrUndefinedColumn)
	}

	if errors.As(err, &validator.ErrValidator) {
		return httphelper.NewHTTPResponse(c, fiber.StatusBadRequest, err)
	}

	log.Println(err.Error())
	return httphelper.NewHTTPResponse(c, fiber.StatusInternalServerError, messages.ErrGeneric)
}

// Creates a new handler.
func NewProfileHandler(route fiber.Router, ps domain.ProfileService, mid *middleware.RequesttMiddleware) {
	handler := &ProfileHandler{
		profileService: ps,
	}

	route.Use(middleware.MidAccess)

	route.Get("", middleware.GetGenericFilter, handler.getProfiles)
	route.Post("", middleware.GetProfileDTO, handler.createProfile)
	route.Get("/:"+httphelper.ParamID, mid.ProfileByID, handler.getProfile)
	route.Put("/:"+httphelper.ParamID, mid.ProfileByID, middleware.GetProfileDTO, handler.updateProfile)
	route.Delete("/:"+httphelper.ParamID, mid.ProfileByID, handler.deleteProfile)
}

// getProfiles godoc
// @Summary      Get profiles
// @Description  Get profiles
// @Tags         Profile
// @Accept       json
// @Produce      json
// @Param        lang query string false "Language responses"
// @Param        filter query filter.Filter false "Optional Filter"
// @Success      200  {array}   dto.ItemsOutputDTO
// @Failure      500  {object}  httphelper.HTTPResponse
// @Router       /profile [get]
// @Security	 Bearer
func (h *ProfileHandler) getProfiles(c *fiber.Ctx) error {
	response, err := h.profileService.GetProfilesOutputDTO(c.Context(), c.Locals(httphelper.LocalFilter).(*filter.Filter))
	if err != nil {
		return h.handlerError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(response)
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
	profileDTO := c.Locals(httphelper.LocalDTO).(*dto.ProfileInputDTO)
	profile, err := h.profileService.CreateProfile(c.Context(), profileDTO)
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
	profileDTO := c.Locals(httphelper.LocalDTO).(*dto.ProfileInputDTO)
	profile := c.Locals(httphelper.LocalObject).(*domain.Profile)
	if err := h.profileService.UpdateProfile(c.Context(), profile, profileDTO); err != nil {
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
