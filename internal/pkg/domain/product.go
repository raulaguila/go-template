package domain

import (
	"context"

	"github.com/raulaguila/go-template/internal/pkg/dto"
	gormhelper "github.com/raulaguila/go-template/pkg/gorm-helper"
	"github.com/raulaguila/go-template/pkg/validator"
)

const ProductTableName string = "product"

type (
	Product struct {
		Base
		Name string `json:"name" gorm:"column:name;type:varchar(100);unique;index;not null;" validate:"required,min=2"`
	}

	ProductRepository interface {
		GetProductByID(context.Context, uint) (*Product, error)
		GetCategories(context.Context, *gormhelper.Filter) (*[]Product, error)
		CountCategories(context.Context, *gormhelper.Filter) (int64, error)
		CreateProduct(context.Context, *dto.ProductInputDTO) (uint, error)
		UpdateProduct(context.Context, *Product, *dto.ProductInputDTO) error
		DeleteProduct(context.Context, *Product) error
	}

	ProductService interface {
		GetProductByID(context.Context, uint) (*Product, error)
		GetCategories(context.Context, *gormhelper.Filter) (*[]Product, error)
		CountCategories(context.Context, *gormhelper.Filter) (int64, error)
		CreateProduct(context.Context, *dto.ProductInputDTO) (uint, error)
		UpdateProduct(context.Context, *Product, *dto.ProductInputDTO) error
		DeleteProduct(context.Context, *Product) error
	}
)

func (Product) TableName() string {
	return ProductTableName
}

func (s *Product) Bind(p *dto.ProductInputDTO) error {
	if p.Name != nil {
		s.Name = *p.Name
	}

	return validator.StructValidator.Validate(s)
}

func (s Product) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"name": s.Name,
	}
}
