package entity

import (
	"time"

	"github.com/mikiasgoitom/Secure-Asset/internal/domain/valueobject"
)

// User represents a user in the system.
type User struct {
	ID              string                     `bson:"_id,omitempty"`
	Username        string                     `bson:"username" validate:"required"`
	Email           string                     `bson:"email" validate:"required,email"`
	Password        string                     `bson:"password" validate:"required"`
	Role            string                     `bson:"role" validate:"required"`
	Department      string                     `bson:"department"`
	ClearanceLevel  valueobject.Classification `bson:"clearanceLevel"` // For MAC (e.g., "Confidential", "Restricted", "Internal", "Public")
	IsMFAEnabled    bool                       `bson:"isMFAEnabled"`
	OTPSecret       string                     `bson:"otpSecret,omitempty"`
	OTPEnabled      bool                       `bson:"otpEnabled"`
	IsAccountLocked bool                       `bson:"isAccountLocked"`
	CreatedAt       time.Time                  `bson:"createdAt"`
	UpdatedAt       time.Time                  `bson:"updatedAt"`
	// Token           string `bson:"token,omitempty"`
}
