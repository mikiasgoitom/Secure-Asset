package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/mikiasgoitom/Secure-Asset/internal/contract"
	"github.com/mikiasgoitom/Secure-Asset/internal/domain/entity"
	"github.com/mikiasgoitom/Secure-Asset/internal/domain/valueobject"
)

type AssetUsecase struct {
	assetRepository contract.IAssetRepository
	userRepository  contract.IUserRepository
	Logger          contract.ILogger
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

// GetAsset retrieves an asset and enforces MAC, ABAC, and RuBAC.
func (u *AssetUsecase) GetAsset(ctx context.Context, assetID string) (*entity.Asset, error) {
	// Retrieve the asset from the database
	asset, err := u.assetRepository.GetByID(ctx, assetID)

	if err != nil {
		u.Logger.Error("Database error fetching asset", valueobject.Field{Key: "assetID", Value: assetID}, valueobject.Field{Key: "error", Value: err})
		return nil, errors.New("internal server error")
	}
	if asset == nil {
		return nil, errors.New("asset not found")
	}

	// Get user's clearance level from the context (set by auth middleware)
	var userClearance valueobject.Classification
	switch v := ctx.Value("clearanceLevel").(type) {
	case valueobject.Classification:
		userClearance = v
	case string:
		userClearance = valueobject.GetClassificationLevel(v)
	default:
		u.Logger.Error("user clearance level not found or invalid type", valueobject.Field{Key: "contextValue", Value: v})
		return nil, errors.New("user clearance level not found in context")
	}

	u.Logger.Debug("check all user contexts", valueobject.Field{Key: "clearanceLevel", Value: userClearance}, valueobject.Field{Key: "role", Value: ctx.Value("role")})


	// --- Mandatory Access Control (MAC) ---
	// User's clearance must be >= asset's classification.
    
	if userClearance < asset.Classification {
        
		u.Logger.Warn("MAC violation attempt",
			valueobject.Field{Key: "assetID", Value: assetID},
			valueobject.Field{Key: "assetClassification", Value: asset.Classification},
			valueobject.Field{Key: "userClearance", Value: userClearance},
		)
		// that the asset exists but the user can't access it.
		return nil, errors.New("asset not found or access denied!!")
	}

	// --- Attribute-Based Access Control (ABAC) ---
	// Only users from the 'IT' department can view 'Server' type assets.
	userDepartment, _ := ctx.Value("department").(string)
    if asset.AssetType == "Server" && userDepartment != "IT" {
        u.Logger.Warn("ABAC violation attempt",
            valueobject.Field{Key: "assetID", Value: assetID},
            valueobject.Field{Key: "reason", Value: "Non-IT user tried to access a server"},
        )
        return nil, errors.New("asset not found")
    }

	// --- Rule-Based Access Control (RuBAC) ---
	// 'Confidential' assets can only be accessed during working hours (9 AM to 5 PM).
	if asset.Classification == valueobject.Confidential {
		now := time.Now()
		if now.Hour() < 9 || now.Hour() >= 17 {
			u.Logger.Warn("RuBAC violation attempt",
				valueobject.Field{Key: "assetID", Value: assetID},
				valueobject.Field{Key: "reason", Value: "Off-hours access to Confidential asset"},
			)
			return nil, errors.New("asset not found")
		}
	}

	u.Logger.Info("Asset retrieved successfully", valueobject.Field{Key: "assetID", Value: assetID})
	return asset, nil
}
