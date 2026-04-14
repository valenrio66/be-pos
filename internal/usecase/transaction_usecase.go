package usecase

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/uptrace/bun"
	"github.com/valenrio66/be-pos/internal/domain"
	"go.uber.org/zap"
)

type TransactionUsecase struct {
	db          *bun.DB
	transRepo   domain.TransactionRepository
	productRepo domain.ProductRepository
	log         *zap.Logger
}

func NewTransactionUsecase(db *bun.DB, transRepo domain.TransactionRepository, productRepo domain.ProductRepository, log *zap.Logger) *TransactionUsecase {
	return &TransactionUsecase{db: db, transRepo: transRepo, productRepo: productRepo, log: log}
}

func (u *TransactionUsecase) CreateTransaction(ctx context.Context, cashierID int64, items []domain.TransactionDetail) (*domain.Transaction, error) {

	tx, err := u.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to start database transaction")
	}
	defer func() {
		errRollback := tx.Rollback()
		if errRollback != nil && !errors.Is(errRollback, sql.ErrTxDone) {
			u.log.Warn("Failed to execute database rollback", zap.Error(errRollback))
		}
	}()

	var totalPrice float64
	transactionNo := fmt.Sprintf("TRX-%d", time.Now().Unix())

	newTrans := &domain.Transaction{
		TransactionNo: transactionNo,
		CashierID:     cashierID,
		TotalPrice:    0,
	}

	if err := u.transRepo.CreateWithTx(ctx, tx, newTrans); err != nil {
		return nil, fmt.Errorf("failed to create transaction header")
	}

	for _, item := range items {
		product, err := u.productRepo.FindByID(ctx, item.ProductID)
		if err != nil {
			return nil, fmt.Errorf("product with ID %d not found", item.ProductID)
		}

		if product.Stock < item.Qty {
			return nil, fmt.Errorf("insufficient stock for product: %s", product.Name)
		}

		subtotal := product.Price * float64(item.Qty)
		totalPrice += subtotal

		item.TransactionID = newTrans.ID
		item.PriceAtBuy = product.Price
		item.Subtotal = subtotal

		if err := u.transRepo.CreateDetailWithTx(ctx, tx, &item); err != nil {
			return nil, fmt.Errorf("failed to save transaction detail")
		}

		product.Stock -= item.Qty
		if err := u.productRepo.Update(ctx, product); err != nil {
			return nil, fmt.Errorf("failed to update product stock")
		}
	}

	newTrans.TotalPrice = totalPrice
	_, err = tx.NewUpdate().Model(newTrans).WherePK().Exec(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to finalize transaction total")
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction")
	}

	return newTrans, nil
}
