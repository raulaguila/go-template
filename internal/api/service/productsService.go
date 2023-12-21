package service

import (
	"context"

	"github.com/raulaguila/go-template/internal/pkg/domain"
	"github.com/raulaguila/go-template/internal/pkg/dto"
	"github.com/raulaguila/go-template/pkg/filter"
)

func NewProductService(r domain.ProductRepository) domain.ProductService {
	return &productService{
		productRepository: r,
	}
}

type productService struct {
	productRepository domain.ProductRepository
}

// Implementation of 'GetProductByID'.
func (s *productService) GetProductByID(ctx context.Context, productID uint) (*domain.Product, error) {
	return s.productRepository.GetProductByID(ctx, productID)
}

// Implementation of 'GetProducts'.
func (s *productService) GetProducts(ctx context.Context, filter *filter.Filter) (*[]domain.Product, error) {
	return s.productRepository.GetProducts(ctx, filter)
}

// Implementation of 'CountProducts'.
func (s *productService) CountProducts(ctx context.Context, filter *filter.Filter) (int64, error) {
	return s.productRepository.CountProducts(ctx, filter)
}

// Implementation of 'CreateProduct'.
func (s *productService) CreateProduct(ctx context.Context, datas *dto.ProductInputDTO) (uint, error) {
	return s.productRepository.CreateProduct(ctx, datas)
}

// Implementation of 'UpdateProduct'.
func (s *productService) UpdateProduct(ctx context.Context, product *domain.Product, datas *dto.ProductInputDTO) error {
	return s.productRepository.UpdateProduct(ctx, product, datas)
}

// Implementation of 'DeleteProduct'.
func (s *productService) DeleteProduct(ctx context.Context, product *domain.Product) error {
	return s.productRepository.DeleteProduct(ctx, product)
}
