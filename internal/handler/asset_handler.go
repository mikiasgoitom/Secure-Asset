package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mikiasgoitom/Secure-Asset/internal/contract"
	"github.com/mikiasgoitom/Secure-Asset/internal/domain/valueobject"
	"github.com/mikiasgoitom/Secure-Asset/internal/dto"
)

type AssetHandler struct {
	assetUsecase contract.IAssetUsecase
	logger       contract.ILogger
}

func NewAssetHandler(assetUsecase contract.IAssetUsecase, logger contract.ILogger) *AssetHandler {
	return &AssetHandler{
		assetUsecase: assetUsecase,
		logger:       logger,
	}
}

func (h *AssetHandler) CreateAsset(ctx *gin.Context) {
	var req dto.CreateAssetRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Failed to bind create asset request", valueobject.Field{Key: "Error", Value: err})
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	ownerID, exist := ctx.Get("userID")
	if !exist {
		h.logger.Error("User ID not found in context")
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	asset, err := h.assetUsecase.CreateAsset(ctx.Request.Context(), req.Name, req.AssetType, req.Classification, req.OwnerUsername)
	if err != nil {
		h.logger.Error("Failed to create asset", valueobject.Field{Key: "Error", Value: err})
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create asset"})
		return
	}
	h.logger.Info("Asset created successfully", valueobject.Field{Key: "AssetID", Value: asset.ID}, valueobject.Field{Key: "OwnerID", Value: ownerID})
	ctx.JSON(http.StatusCreated, gin.H{"asset": asset})
}