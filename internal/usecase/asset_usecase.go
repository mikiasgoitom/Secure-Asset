package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/mikiasgoitom/Secure-Asset/internal/contract"
	"github.com/mikiasgoitom/Secure-Asset/internal/domain/entity"
	"github.com/mikiasgoitom/Secure-Asset/internal/domain/valueobject"
)

type AssetUsecase struct {
	assetRepository contract.IAssetRepository
	userRepository  contract.IUserRepository
	Logger		  contract.ILogger
}

func NewAssetUsecase(assetRepo contract.IAssetRepository, userRepo contract.IUserRepository, logger contract.ILogger) contract.IAssetUsecase {
	return &AssetUsecase{
		assetRepository: assetRepo,
		userRepository:  userRepo,
		Logger:          logger,
	}
}

func (uc *AssetUsecase) CreateAsset(ctx context.Context, name, assetType string, classification uint8, ownerID string) (*entity.Asset, error) {
	_, err := uc.userRepository.FindByUsername(ctx, ownerID)
	if err != nil {
		uc.Logger.Error("Failed to find user by username", valueobject.Field{Key: "error", Value: err})
		return nil, err
	}
	if classification > uint8(valueobject.Public) || classification < uint8(valueobject.Confidential) {
		uc.Logger.Error("Invalid classification level", valueobject.Field{Key: "classification", Value: classification})
		return nil, err
	}
	newAsset := &entity.Asset{
		ID:             uuid.New().String(),
		Name:           name,
		AssetType:      assetType,
		Classification: valueobject.Classification(classification),
		OwnerID:        ownerID,
	}
	createdAsset, err := uc.assetRepository.Create(ctx, newAsset)
	if err != nil {
		uc.Logger.Error("Failed to create asset", valueobject.Field{Key: "error", Value: err})
		return nil, err
	}
	return createdAsset, nil
}