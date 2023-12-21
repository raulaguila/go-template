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
	"github.com/raulaguila/go-template/pkg/postgresql"
	"github.com/raulaguila/go-template/pkg/validator"
)

type ProductHandler struct {
	productService domain.ProductService
}

func (ProductHandler) foreignKeyViolatedMethod(c *fiber.Ctx, translation *i18n.Translation) error {
	switch c.Method() {
	case fiber.MethodPut, fiber.MethodPost, fiber.MethodPatch:
		return httphelper.NewHTTPErrorResponse(c, fiber.StatusBadRequest, translation.ErrProductNotFound)
	case fiber.MethodDelete:
		return httphelper.NewHTTPErrorResponse(c, fiber.StatusBadRequest, translation.ErrProductUsed)
	default:
		return httphelper.NewHTTPErrorResponse(c, fiber.StatusInternalServerError, translation.ErrGeneric)
	}
}

func (s ProductHandler) handlerError(c *fiber.Ctx, err error) error {
	translation := c.Locals(httphelper.LocalLang).(*i18n.Translation)

	switch postgresql.HandlerError(err) {
	case postgresql.ErrDuplicatedKey:
		return httphelper.NewHTTPErrorResponse(c, fiber.StatusConflict, translation.ErrProductRegistered)
	case postgresql.ErrForeignKeyViolated:
		return s.foreignKeyViolatedMethod(c, translation)
	case postgresql.ErrUndefinedColumn:
		return httphelper.NewHTTPErrorResponse(c, fiber.StatusBadRequest, translation.ErrUndefinedColumn)
	}

	if errors.As(err, &validator.ErrValidator) {
		return httphelper.NewHTTPErrorResponse(c, fiber.StatusBadRequest, err)
	}

	log.Println(err.Error())
	return httphelper.NewHTTPErrorResponse(c, fiber.StatusInternalServerError, translation.ErrGeneric)
}

// Creates a new handler.
func NewProductHandler(route fiber.Router, ps domain.ProductService, mid *middleware.RequesttMiddleware) {
	handler := &ProductHandler{
		productService: ps,
	}

	route.Use(middleware.MidAccess)
	productByID := mid.ItemByID(&domain.Product{}, domain.ProductTableName)

	route.Get("", middleware.GetGenericFilter, handler.getProducts)
	route.Post("", middleware.GetDTO(&dto.ProductInputDTO{}), handler.createProduct)
	route.Get("/:"+httphelper.ParamID, productByID, handler.getProductBydID)
	route.Put("/:"+httphelper.ParamID, productByID, middleware.GetDTO(&dto.ProductInputDTO{}), handler.updateProduct)
	route.Delete("/:"+httphelper.ParamID, productByID, handler.deleteProduct)
}

// getProducts godoc
// @Summary      Get products
// @Description  Get products
// @Tags         Product
// @Accept       json
// @Produce      json
// @Param        lang query string false "Language responses"
// @Param        filter query filter.Filter false "Optional Filter"
// @Success      200  {array}   dto.ItemsOutputDTO
// @Failure      500  {object}  httphelper.HTTPResponse
// @Router       /product [get]
// @Security	 Bearer
func (h *ProductHandler) getProducts(c *fiber.Ctx) error {
	products, err := h.productService.GetProducts(c.Context(), c.Locals(httphelper.LocalFilter).(*filter.Filter))
	if err != nil {
		return h.handlerError(c, err)
	}

	count, err := h.productService.CountProducts(c.Context(), c.Locals(httphelper.LocalFilter).(*filter.Filter))
	if err != nil {
		return h.handlerError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(&dto.ItemsOutputDTO{
		Items: products,
		Count: count,
	})
}

// getProductBydID godoc
// @Summary      Get product by ID
// @Description  Get product by ID
// @Tags         Product
// @Accept       json
// @Produce      json
// @Param        lang query string false "Language responses"
// @Param        id   path			int			true        "Product ID"
// @Success      200  {object}  domain.Product
// @Failure      400  {object}  httphelper.HTTPResponse
// @Failure      404  {object}  httphelper.HTTPResponse
// @Failure      500  {object}  httphelper.HTTPResponse
// @Router       /product/{id} [get]
// @Security	 Bearer
func (h *ProductHandler) getProductBydID(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(c.Locals(httphelper.LocalObject).(*domain.Product))
}

// createProduct godoc
// @Summary      Insert product
// @Description  Insert product
// @Tags         Product
// @Accept       json
// @Produce      json
// @Param        lang query string false "Language responses"
// @Param        product body dto.ProductInputDTO true "Product model"
// @Success      201  {object}  domain.Product
// @Failure      400  {object}  httphelper.HTTPResponse
// @Failure      409  {object}  httphelper.HTTPResponse
// @Failure      500  {object}  httphelper.HTTPResponse
// @Router       /product [post]
// @Security	 Bearer
func (h *ProductHandler) createProduct(c *fiber.Ctx) error {
	id, err := h.productService.CreateProduct(c.Context(), c.Locals(httphelper.LocalDTO).(*dto.ProductInputDTO))
	if err != nil {
		return h.handlerError(c, err)
	}

	product, err := h.productService.GetProductByID(c.Context(), id)
	if err != nil {
		return h.handlerError(c, err)
	}

	return c.Status(fiber.StatusCreated).JSON(product)
}

// updateProduct godoc
// @Summary      Update product by ID
// @Description  Update product by ID
// @Tags         Product
// @Accept       json
// @Produce      json
// @Param        lang query string false "Language responses"
// @Param        id     path    int     true        "Product ID"
// @Param        product body dto.ProductInputDTO true "Product model"
// @Success      200  {object}  domain.Product
// @Failure      400  {object}  httphelper.HTTPResponse
// @Failure      404  {object}  httphelper.HTTPResponse
// @Failure      500  {object}  httphelper.HTTPResponse
// @Router       /product/{id} [put]
// @Security	 Bearer
func (h *ProductHandler) updateProduct(c *fiber.Ctx) error {
	product := c.Locals(httphelper.LocalObject).(*domain.Product)
	if err := h.productService.UpdateProduct(c.Context(), product, c.Locals(httphelper.LocalDTO).(*dto.ProductInputDTO)); err != nil {
		return h.handlerError(c, err)
	}

	updated, err := h.productService.GetProductByID(c.Context(), product.Id)
	if err != nil {
		return h.handlerError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(updated)
}

// deleteProduct godoc
// @Summary      Delete product by ID
// @Description  Delete product by ID
// @Tags         Product
// @Accept       json
// @Produce      json
// @Param        lang query string false "Language responses"
// @Param        id     path    int     true        "Product ID"
// @Success      204  {object}  nil
// @Failure      404  {object}  httphelper.HTTPResponse
// @Failure      500  {object}  httphelper.HTTPResponse
// @Router       /product/{id} [delete]
// @Security	 Bearer
func (h *ProductHandler) deleteProduct(c *fiber.Ctx) error {
	if err := h.productService.DeleteProduct(c.Context(), c.Locals(httphelper.LocalObject).(*domain.Product)); err != nil {
		return h.handlerError(c, err)
	}

	return c.Status(fiber.StatusNoContent).Send(nil)
}
