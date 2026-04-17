package dto

type TransactionItemRequest struct {
	ProductID int64 `json:"product_id" binding:"required"`
	Qty       int   `json:"qty" binding:"required,gt=0"`
}

type CreateTransactionRequest struct {
	Items []TransactionItemRequest `json:"items" binding:"required,min=1"`
}

type TransactionItemResponse struct {
	ProductID  int64   `json:"product_id"`
	Qty        int     `json:"qty"`
	PriceAtBuy float64 `json:"price_at_buy"`
	Subtotal   float64 `json:"subtotal"`
}

type TransactionResponse struct {
	ID            int64                     `json:"id"`
	TransactionNo string                    `json:"transaction_no"`
	TotalPrice    float64                   `json:"total_price"`
	CashierID     int64                     `json:"cashier_id"`
	CreatedAt     string                    `json:"created_at"`
	Items         []TransactionItemResponse `json:"items,omitempty"`
}

type DashboardSummaryResponse struct {
	Date              string  `json:"date"`
	TotalTransactions int     `json:"total_transactions"`
	TotalRevenue      float64 `json:"total_revenue"`
}
