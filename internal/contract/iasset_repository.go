package contract

import (
	"context"

	"github.com/mikiasgoitom/Secure-Asset/internal/domain/entity"
)

type IAssetRepository interface {
	Create(ctx context.Context, asset *entity.Asset) (*entity.Asset, error)
	GetByID(ctx context.Context, id string) (*entity.Asset, error)
	Update(ctx context.Context, asset *entity.Asset) (*entity.Asset, error)
	// List(ctx context.Context, filter map[string]interface{}) ([]*entity.Asset, error)
}