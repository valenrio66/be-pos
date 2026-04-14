package postgres

import (
	"context"

	"github.com/uptrace/bun"
	"github.com/valenrio66/be-pos/internal/domain"
)

type transactionRepositoryImpl struct {
	db *bun.DB
}

func NewTransactionRepository(db *bun.DB) domain.TransactionRepository {
	return &transactionRepositoryImpl{db: db}
}

func (r *transactionRepositoryImpl) CreateWithTx(ctx context.Context, tx bun.Tx, trans *domain.Transaction) error {
	_, err := tx.NewInsert().Model(trans).Exec(ctx)
	return err
}

func (r *transactionRepositoryImpl) CreateDetailWithTx(ctx context.Context, tx bun.Tx, detail *domain.TransactionDetail) error {
	_, err := tx.NewInsert().Model(detail).Exec(ctx)
	return err
}

func (r *transactionRepositoryImpl) FindAll(ctx context.Context) ([]domain.Transaction, error) {
	var list []domain.Transaction
	err := r.db.NewSelect().Model(&list).Order("created_at DESC").Scan(ctx)
	return list, err
}

func (r *transactionRepositoryImpl) FindByID(ctx context.Context, id int64) (*domain.Transaction, error) {
	trans := new(domain.Transaction)
	err := r.db.NewSelect().Model(trans).Relation("Items").Where("t.id = ?", id).Scan(ctx)
	return trans, err
}
