package entity

import "github.com/mikiasgoitom/Secure-Asset/internal/domain/valueobject"

// Asset represents a company asset.
type Asset struct {
	ID               string                     `bson:"_id,omitempty"`
	Name             string                     `bson:"name" validate:"required"`
	AssetType        string                     `bson:"assetType" validate:"required"` // E.g., "Laptop", "Server", "Software License"
	OwnerID          string                     `bson:"ownerId"`                       // For DAC - The user who owns/manages this asset
	Department       string                     `bson:"department"`
	Classification   valueobject.Classification `bson:"classification" validate:"required"` // For MAC (e.g., "Confidential", "Restricted", "Internal", "Public")
	CurrentStatus    string                     `bson:"currentStatus"`                      // E.g., "Available", "In Use", "Under Maintenance"
	CurrentUserID    string                     `bson:"currentUserId,omitempty"`            // Who is currently using the asset
	RequestHistory   []AssetRequest             `bson:"requestHistory,omitempty"`
	Location         string                     `bson:"location"`                   // For RuBAC/ABAC
	Value            float64                    `bson:"value"`                      // For RuBAC (e.g., high-value assets)
	CustomAttributes map[string]string          `bson:"customAttributes,omitempty"` // For ABAC
}

// AssetRequest represents a request to borrow an asset.
type AssetRequest struct {
	RequestorID string `bson:"requestorId"`
	RequestDate string `bson:"requestDate"`
	ReturnDate  string `bson:"returnDate,omitempty"`
	Status      string `bson:"status"` // "Pending", "Approved", "Rejected", "Returned"
	ApproverID  string `bson:"approverId,omitempty"`
}
