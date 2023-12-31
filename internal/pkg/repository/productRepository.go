package repository

import (
	"context"

	"github.com/raulaguila/go-template/internal/pkg/domain"
	"github.com/raulaguila/go-template/internal/pkg/dto"
	"github.com/raulaguila/go-template/pkg/filter"
	"gorm.io/gorm"
)

func NewProductRepository(postgres *gorm.DB) domain.ProductRepository {
	return &productRepository{
		postgres: postgres,
	}
}

type productRepository struct {
	postgres *gorm.DB
}

func (s *productRepository) applyFilter(ctx context.Context, filter *filter.Filter) *gorm.DB {
	postgres := s.postgres.WithContext(ctx)
	postgres = filter.ApplySearchLike(postgres, "name")
	postgres = filter.ApplyOrder(postgres)

	return postgres
}

func (s *productRepository) countProducts(postgres *gorm.DB) (int64, error) {
	var count int64
	return count, postgres.Model(&domain.Product{}).Count(&count).Error
}

func (s *productRepository) listProducts(postgres *gorm.DB) (*[]domain.Product, error) {
	products := &[]domain.Product{}
	return products, postgres.Find(products).Error
}

func (s *productRepository) GetProductsOutputDTO(ctx context.Context, filter *filter.Filter) (*dto.ItemsOutputDTO, error) {
	postgres := s.applyFilter(ctx, filter)
	count, err := s.countProducts(postgres)
	if err != nil {
		return nil, err
	}

	postgres = filter.ApplyPagination(postgres)
	items, err := s.listProducts(postgres)
	if err != nil {
		return nil, err
	}

	return &dto.ItemsOutputDTO{
		Items: items,
		Count: count,
	}, nil
}

func (s *productRepository) GetProductByID(ctx context.Context, productID uint) (*domain.Product, error) {
	product := &domain.Product{}
	return product, s.postgres.WithContext(ctx).First(product, productID).Error
}

func (s *productRepository) CreateProduct(ctx context.Context, datas *dto.ProductInputDTO) (*domain.Product, error) {
	product := &domain.Product{}
	if err := product.Bind(datas); err != nil {
		return nil, err
	}

	return product, s.postgres.WithContext(ctx).Create(product).Error
}

func (s *productRepository) UpdateProduct(ctx context.Context, product *domain.Product, datas *dto.ProductInputDTO) error {
	if err := product.Bind(datas); err != nil {
		return err
	}

	return s.postgres.WithContext(ctx).Model(product).Updates(product.ToMap()).Error
}

func (s *productRepository) DeleteProduct(ctx context.Context, product *domain.Product) error {
	return s.postgres.WithContext(ctx).Delete(product).Error
}
