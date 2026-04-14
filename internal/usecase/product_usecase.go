package usecase

import (
	"context"
	"database/sql"
	"errors"

	"github.com/valenrio66/be-pos/internal/domain"
	"go.uber.org/zap"
)

type ProductUsecase struct {
	repo domain.ProductRepository
	log  *zap.Logger
}

func NewProductUsecase(repo domain.ProductRepository, log *zap.Logger) *ProductUsecase {
	return &ProductUsecase{repo: repo, log: log}
}

func (u *ProductUsecase) CreateProduct(ctx context.Context, product *domain.Product) error {
	err := u.repo.Create(ctx, product)
	if err != nil {
		u.log.Error("Failed to create product", zap.Error(err))
		return errors.New("failed to save product to database")
	}
	return nil
}

func (u *ProductUsecase) GetAllProducts(ctx context.Context) ([]domain.Product, error) {
	products, err := u.repo.FindAll(ctx)
	if err != nil {
		u.log.Error("Failed to fetch products", zap.Error(err))
		return nil, errors.New("failed to retrieve products")
	}
	return products, nil
}

func (u *ProductUsecase) GetProductByID(ctx context.Context, id int64) (*domain.Product, error) {
	product, err := u.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("product not found")
		}
		u.log.Error("Failed to fetch product by ID", zap.Error(err))
		return nil, errors.New("failed to retrieve product")
	}
	return product, nil
}

func (u *ProductUsecase) UpdateProduct(ctx context.Context, id int64, updatedData *domain.Product) error {
	// 1. Cek apakah produk ada
	existingProduct, err := u.GetProductByID(ctx, id)
	if err != nil {
		return err
	}

	// 2. Timpa data lama dengan data baru
	existingProduct.SKU = updatedData.SKU
	existingProduct.Name = updatedData.Name
	existingProduct.Description = updatedData.Description
	existingProduct.Price = updatedData.Price
	existingProduct.Unit = updatedData.Unit
	existingProduct.Stock = updatedData.Stock

	// 3. Simpan perubahan
	err = u.repo.Update(ctx, existingProduct)
	if err != nil {
		u.log.Error("Failed to update product", zap.Error(err))
		return errors.New("failed to update product in database")
	}
	return nil
}

func (u *ProductUsecase) DeleteProduct(ctx context.Context, id int64) error {
	// Cek apakah produk ada sebelum dihapus
	_, err := u.GetProductByID(ctx, id)
	if err != nil {
		return err
	}

	err = u.repo.Delete(ctx, id)
	if err != nil {
		u.log.Error("Failed to delete product", zap.Error(err))
		return errors.New("failed to delete product")
	}
	return nil
}
