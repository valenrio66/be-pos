package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/valenrio66/be-pos/internal/delivery/http/dto"
	"github.com/valenrio66/be-pos/internal/domain"
	"github.com/valenrio66/be-pos/internal/usecase"
	"github.com/valenrio66/be-pos/pkg/response"
	"github.com/valenrio66/be-pos/pkg/utils"
)

type TransactionHandler struct {
	usecase *usecase.TransactionUsecase
}

func NewTransactionHandler(usecase *usecase.TransactionUsecase) *TransactionHandler {
	return &TransactionHandler{usecase: usecase}
}

func (h *TransactionHandler) Checkout(c *gin.Context) {
	var req dto.CreateTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid input data", utils.FormatValidationError(err))
		return
	}

	cashierID := c.MustGet("user_id").(int64)

	var domainItems []domain.TransactionDetail
	for _, item := range req.Items {
		domainItems = append(domainItems, domain.TransactionDetail{
			ProductID: item.ProductID,
			Qty:       item.Qty,
		})
	}

	res, err := h.usecase.CreateTransaction(c.Request.Context(), cashierID, domainItems)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Transaction failed", err.Error())
		return
	}

	resp := dto.TransactionResponse{
		ID:            res.ID,
		TransactionNo: res.TransactionNo,
		TotalPrice:    res.TotalPrice,
		CashierID:     res.CashierID,
		CreatedAt:     res.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	response.Success(c, http.StatusCreated, "Transaction completed successfully", resp)
}
