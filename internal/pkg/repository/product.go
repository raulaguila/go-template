package repository

import (
	"context"

	"github.com/raulaguila/go-template/internal/pkg/domain"
	"github.com/raulaguila/go-template/internal/pkg/dto"
	gormhelper "github.com/raulaguila/go-template/pkg/gorm-helper"
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

func (s *productRepository) applyFilter(ctx context.Context, filter *gormhelper.Filter, pag bool) *gorm.DB {
	postgres := s.postgres.WithContext(ctx)
	postgres = filter.ApplySearchLike(postgres, "name")
	postgres = filter.ApplyOrder(postgres)
	if pag {
		postgres = filter.ApplyPagination(postgres)
	}

	return postgres
}

func (s *productRepository) GetProductByID(ctx context.Context, productID uint) (*domain.Product, error) {
	product := &domain.Product{}
	return product, s.postgres.WithContext(ctx).First(product, productID).Error
}

func (s *productRepository) GetProducts(ctx context.Context, filter *gormhelper.Filter) (*[]domain.Product, error) {
	products := &[]domain.Product{}
	return products, s.applyFilter(ctx, filter, true).Find(products).Error
}

func (s *productRepository) CountProducts(ctx context.Context, filter *gormhelper.Filter) (int64, error) {
	var count int64
	return count, s.applyFilter(ctx, filter, false).Model(&domain.Product{}).Count(&count).Error
}

func (s *productRepository) CreateProduct(ctx context.Context, datas *dto.ProductInputDTO) (uint, error) {
	product := &domain.Product{}
	if err := product.Bind(datas); err != nil {
		return 0, err
	}
	return product.Id, s.postgres.WithContext(ctx).Create(product).Error
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
