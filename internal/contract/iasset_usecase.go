package contract

import (
	"context"

	"github.com/mikiasgoitom/Secure-Asset/internal/domain/entity"
)
	
type IAssetUsecase interface {
	CreateAsset(ctx context.Context, name, assetType string, classification uint8, ownerID string) (*entity.Asset, error)
}