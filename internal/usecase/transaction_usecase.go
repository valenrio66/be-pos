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

func (u *TransactionUsecase) CreateTransaction(ctx context.Context, cashierID int64, paidAmount float64, items []domain.TransactionDetail) (*domain.Transaction, error) {
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
		PaidAmount:    paidAmount,
		ChangeAmount:  0,
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

	if paidAmount < totalPrice {
		return nil, fmt.Errorf("insufficient payment: total is %.2f but paid amount is %.2f", totalPrice, paidAmount)
	}

	changeAmount := paidAmount - totalPrice

	newTrans.TotalPrice = totalPrice
	newTrans.ChangeAmount = changeAmount

	_, err = tx.NewUpdate().Model(newTrans).WherePK().Exec(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to finalize transaction total")
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction")
	}

	return newTrans, nil
}

func (u *TransactionUsecase) GetTodaySummary(ctx context.Context) (*domain.DailySummary, error) {
	now := time.Now()

	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	endOfDay := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 999999999, now.Location())

	summary, err := u.transRepo.GetDailySummary(ctx, startOfDay, endOfDay)
	if err != nil {
		u.log.Error("Failed to fetch daily summary", zap.Error(err))
		return nil, fmt.Errorf("failed to retrieve dashboard summary")
	}

	return summary, nil
}

func (u *TransactionUsecase) Inquiry(ctx context.Context, items []domain.TransactionDetail) (*domain.TransactionInquiry, error) {
	var totalPrice float64
	var detailedItems []domain.TransactionDetail

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

		detailedItems = append(detailedItems, domain.TransactionDetail{
			ProductID:  product.ID,
			Qty:        item.Qty,
			PriceAtBuy: product.Price,
			Subtotal:   subtotal,
		})
	}

	return &domain.TransactionInquiry{
		TotalPrice: totalPrice,
		Items:      detailedItems,
	}, nil
}
