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
		response.Error(c, http.StatusBadRequest, "Validasi input gagal", formattedErrors)
		return
	}

	user, err := h.usecase.CreateUser(c.Request.Context(), req.Name, req.Email, req.Password)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Gagal membuat user", err.Error())
		return
	}

	response.Success(c, http.StatusCreated, "User berhasil didaftarkan", user)
}

func (h *UserHandler) Login(c *gin.Context) {
	var req dto.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		formattedErrors := utils.FormatValidationError(err)
		response.Error(c, http.StatusBadRequest, "Validasi input gagal", formattedErrors)
		return
	}

	loginResp, err := h.usecase.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		response.Error(c, http.StatusUnauthorized, err.Error(), nil)
		return
	}

	response.Success(c, http.StatusOK, "Login berhasil", loginResp)
}
