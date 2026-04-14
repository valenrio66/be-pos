package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/valenrio66/be-pos/internal/delivery/http/dto"
	"github.com/valenrio66/be-pos/internal/domain"
	"github.com/valenrio66/be-pos/internal/usecase"
	"github.com/valenrio66/be-pos/pkg/response"
	"github.com/valenrio66/be-pos/pkg/utils"
)

type ProductHandler struct {
	usecase *usecase.ProductUsecase
}

func NewProductHandler(usecase *usecase.ProductUsecase) *ProductHandler {
	return &ProductHandler{usecase: usecase}
}

// Helper untuk mengubah string ID dari URL ke int64
func parseID(c *gin.Context) (int64, error) {
	idStr := c.Param("id")
	return strconv.ParseInt(idStr, 10, 64)
}

func (h *ProductHandler) CreateProduct(c *gin.Context) {
	var req dto.ProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid input data", utils.FormatValidationError(err))
		return
	}

	product := &domain.Product{
		SKU:         req.SKU,
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Unit:        req.Unit,
		Stock:       req.Stock,
	}

	if err := h.usecase.CreateProduct(c.Request.Context(), product); err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to process request", err.Error())
		return
	}

	response.Success(c, http.StatusCreated, "Product created successfully", nil)
}

func (h *ProductHandler) GetAllProducts(c *gin.Context) {
	products, err := h.usecase.GetAllProducts(c.Request.Context())
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to process request", err.Error())
		return
	}

	// Mapping ke DTO Response (agar field CreatedAt/UpdatedAt tidak ikut terkirim jika tidak perlu)
	var resp []dto.ProductResponse
	for _, p := range products {
		resp = append(resp, dto.ProductResponse{
			ID:          p.ID,
			SKU:         p.SKU,
			Name:        p.Name,
			Description: p.Description,
			Price:       p.Price,
			Unit:        p.Unit,
			Stock:       p.Stock,
		})
	}

	// Jika kosong, kembalikan array kosong, bukan null
	if resp == nil {
		resp = []dto.ProductResponse{}
	}

	response.Success(c, http.StatusOK, "Products retrieved successfully", resp)
}

func (h *ProductHandler) GetProductByID(c *gin.Context) {
	id, err := parseID(c)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid product ID", nil)
		return
	}

	p, err := h.usecase.GetProductByID(c.Request.Context(), id)
	if err != nil {
		response.Error(c, http.StatusNotFound, err.Error(), nil)
		return
	}

	resp := dto.ProductResponse{
		ID:          p.ID,
		SKU:         p.SKU,
		Name:        p.Name,
		Description: p.Description,
		Price:       p.Price,
		Unit:        p.Unit,
		Stock:       p.Stock,
	}

	response.Success(c, http.StatusOK, "Product retrieved successfully", resp)
}

func (h *ProductHandler) UpdateProduct(c *gin.Context) {
	id, err := parseID(c)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid product ID", nil)
		return
	}

	var req dto.ProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid input data", utils.FormatValidationError(err))
		return
	}

	productData := &domain.Product{
		SKU:         req.SKU,
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Unit:        req.Unit,
		Stock:       req.Stock,
	}

	if err := h.usecase.UpdateProduct(c.Request.Context(), id, productData); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	response.Success(c, http.StatusOK, "Product updated successfully", nil)
}

func (h *ProductHandler) DeleteProduct(c *gin.Context) {
	id, err := parseID(c)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid product ID", nil)
		return
	}

	if err := h.usecase.DeleteProduct(c.Request.Context(), id); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	response.Success(c, http.StatusOK, "Product deleted successfully", nil)
}
