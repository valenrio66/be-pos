package domain

import (
	"context"
	"time"

	"github.com/uptrace/bun"
)

type Transaction struct {
	bun.BaseModel `bun:"table:transactions,alias:t"`

	ID            int64     `bun:"id,pk,autoincrement"`
	TransactionNo string    `bun:"transaction_no,unique,notnull"`
	TotalPrice    float64   `bun:"total_price,notnull"`
	CashierID     int64     `bun:"cashier_id,notnull"`
	CreatedAt     time.Time `bun:"created_at,nullzero,notnull,default:current_timestamp"`

	// Relation
	Items []TransactionDetail `bun:"rel:has-many,join:id=transaction_id"`
}

type TransactionDetail struct {
	bun.BaseModel `bun:"table:transaction_details,alias:td"`

	ID            int64   `bun:"id,pk,autoincrement"`
	TransactionID int64   `bun:"transaction_id,notnull"`
	ProductID     int64   `bun:"product_id,notnull"`
	Qty           int     `bun:"qty,notnull"`
	PriceAtBuy    float64 `bun:"price_at_buy,notnull"`
	Subtotal      float64 `bun:"subtotal,notnull"`
}

type DailySummary struct {
	TotalTransactions int     `bun:"total_transactions"`
	TotalRevenue      float64 `bun:"total_revenue"`
}

type TransactionRepository interface {
	CreateWithTx(ctx context.Context, tx bun.Tx, trans *Transaction) error
	CreateDetailWithTx(ctx context.Context, tx bun.Tx, detail *TransactionDetail) error
	FindAll(ctx context.Context) ([]Transaction, error)
	FindByID(ctx context.Context, id int64) (*Transaction, error)
	GetDailySummary(ctx context.Context, startOfDay, endOfDay time.Time) (*DailySummary, error)
}
