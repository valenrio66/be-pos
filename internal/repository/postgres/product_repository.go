package postgres

import (
	"context"

	"github.com/uptrace/bun"
	"github.com/valenrio66/be-pos/internal/domain"
)

type productRepositoryImpl struct {
	db *bun.DB
}

func NewProductRepository(db *bun.DB) domain.ProductRepository {
	return &productRepositoryImpl{db: db}
}

func (r *productRepositoryImpl) Create(ctx context.Context, product *domain.Product) error {
	_, err := r.db.NewInsert().Model(product).Exec(ctx)
	return err
}

func (r *productRepositoryImpl) FindAll(ctx context.Context) ([]domain.Product, error) {
	var products []domain.Product
	err := r.db.NewSelect().Model(&products).Order("name ASC").Scan(ctx)
	return products, err
}

func (r *productRepositoryImpl) FindByID(ctx context.Context, id int64) (*domain.Product, error) {
	product := new(domain.Product)
	err := r.db.NewSelect().Model(product).Where("id = ?", id).Scan(ctx)
	return product, err
}

func (r *productRepositoryImpl) Update(ctx context.Context, product *domain.Product) error {
	_, err := r.db.NewUpdate().Model(product).WherePK().Exec(ctx)
	return err
}

func (r *productRepositoryImpl) Delete(ctx context.Context, id int64) error {
	_, err := r.db.NewDelete().Model((*domain.Product)(nil)).Where("id = ?", id).Exec(ctx)
	return err
}
