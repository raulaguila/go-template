package httphelper

import "github.com/gofiber/fiber/v2"

type HTTPError struct {
	Code    int    `json:"code" example:"400"`
	Message string `json:"message" example:"status bad request"`
}

func NewHTTPError(c *fiber.Ctx, status int, err error) error {
	return c.Status(status).JSON(&HTTPError{
		Code:    status,
		Message: err.Error(),
	})
}
