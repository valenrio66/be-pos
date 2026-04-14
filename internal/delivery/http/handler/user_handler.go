package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/valenrio66/be-pos/internal/delivery/http/dto"
	"github.com/valenrio66/be-pos/internal/usecase"
	"github.com/valenrio66/be-pos/pkg/response"
	"github.com/valenrio66/be-pos/pkg/utils"
)

type UserHandler struct {
	usecase *usecase.UserUsecase
}

func NewUserHandler(usecase *usecase.UserUsecase) *UserHandler {
	return &UserHandler{usecase: usecase}
}

func (h *UserHandler) CreateUser(c *gin.Context) {
	var req dto.CreateUserRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		formattedErrors := utils.FormatValidationError(err)
		response.Error(c, http.StatusBadRequest, "Input validation failed", formattedErrors)
		return
	}

	user, err := h.usecase.CreateUser(c.Request.Context(), req.Name, req.Email, req.Password)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to create user", err.Error())
		return
	}

	response.Success(c, http.StatusCreated, "The user has successfully registered", user)
}

func (h *UserHandler) Login(c *gin.Context) {
	var req dto.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		formattedErrors := utils.FormatValidationError(err)
		response.Error(c, http.StatusBadRequest, "Input validation failed", formattedErrors)
		return
	}

	loginResp, err := h.usecase.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		response.Error(c, http.StatusUnauthorized, err.Error(), nil)
		return
	}

	response.Success(c, http.StatusOK, "Login successful", loginResp)
}

func (h *UserHandler) RefreshToken(c *gin.Context) {
	var req dto.RefreshTokenRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		formattedErrors := utils.FormatValidationError(err)
		response.Error(c, http.StatusBadRequest, "Input validation failed", formattedErrors)
		return
	}

	resp, err := h.usecase.RefreshAccessToken(c.Request.Context(), req.RefreshToken)
	if err != nil {
		response.Error(c, http.StatusUnauthorized, err.Error(), nil)
		return
	}

	response.Success(c, http.StatusOK, "The token has been successfully updated", resp)
}

func (h *UserHandler) GetProfile(c *gin.Context) {
	userIDVal, exists := c.Get("user_id")
	if !exists {
		response.Error(c, http.StatusInternalServerError, "System error", "User ID not found in context")
		return
	}

	userID := userIDVal.(int64)

	user, err := h.usecase.GetProfile(c.Request.Context(), userID)
	if err != nil {
		response.Error(c, http.StatusNotFound, "Data not found", err.Error())
		return
	}

	resp := dto.UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Role:      user.Role,
		CreatedAt: user.CreatedAt.Format("2006-01-02 15:04:05"),
	}

	response.Success(c, http.StatusOK, "Success get profile", resp)
}

func (h *UserHandler) AdminDashboard(c *gin.Context) {

	data := gin.H{
		"system":  "Building Material POS",
		"status":  "Running smoothly",
		"message": "Welcome to the executive dashboard",
	}

	response.Success(c, http.StatusOK, "Admin dashboard retrieved successfully", data)
}
