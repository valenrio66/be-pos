package domain

import (
	"context"
	"time"

	"github.com/uptrace/bun"
)

type Product struct {
	bun.BaseModel `bun:"table:products,alias:p"`

	ID          int64     `bun:"id,pk,autoincrement"`
	SKU         string    `bun:"sku,unique,notnull"`
	Name        string    `bun:"name,notnull"`
	Description string    `bun:"description"`
	Price       float64   `bun:"price,notnull"`
	Unit        string    `bun:"unit,notnull"`
	Stock       int       `bun:"stock,notnull,default:0"`
	CreatedAt   time.Time `bun:"created_at,nullzero,notnull,default:current_timestamp"`
	UpdatedAt   time.Time `bun:"updated_at,nullzero,notnull,default:current_timestamp"`
}

type ProductRepository interface {
	Create(ctx context.Context, product *Product) error
	FindAll(ctx context.Context) ([]Product, error)
	FindByID(ctx context.Context, id int64) (*Product, error)
	Update(ctx context.Context, product *Product) error
	Delete(ctx context.Context, id int64) error
}
