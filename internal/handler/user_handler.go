package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mikiasgoitom/Secure-Asset/internal/contract"
	"github.com/mikiasgoitom/Secure-Asset/internal/domain/valueobject"
	"github.com/mikiasgoitom/Secure-Asset/internal/dto"
)

type UserHandler struct {
	userUsecase contract.IUserUsecase
	Logger      contract.ILogger
}

func NewUserHandler(userUsecase contract.IUserUsecase, logger contract.ILogger) *UserHandler {
	return &UserHandler{
		userUsecase: userUsecase,
		Logger:      logger,
	}
}

func (h *UserHandler) RegisterUser(ctx *gin.Context) {
	var req dto.RegisterRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		h.Logger.Error("Invalid registration request", valueobject.Field{Key: "error", Value: err})
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	_, err := h.userUsecase.Register(ctx.Request.Context(), req.Username, req.Email, req.Password)
	if err != nil {
		h.Logger.Error("User registration failed", valueobject.Field{Key: "error", Value: err})
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "registration failed"})
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{"message": "user registered successfully"})
}
func (h *UserHandler) LoginUser(ctx *gin.Context) {
	var req dto.LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		h.Logger.Error("Invalid login request", valueobject.Field{Key: "error", Value: err})
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	token, err := h.userUsecase.Login(ctx.Request.Context(), req.Identifier, req.Password)
	if err != nil {
		h.Logger.Error("User login failed", valueobject.Field{Key: "error", Value: err})
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"token": token})
}
