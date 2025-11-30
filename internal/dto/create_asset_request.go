package dto

type CreateAssetRequest struct {
	Name           string `json:"name" binding:"required"`
	AssetType      string `json:"asset_type" binding:"required"`
	Classification uint8  `json:"classification" binding:"required,oneof=0 1 2 3"`
	OwnerUsername        string `json:"owner_username" binding:"required"`
}