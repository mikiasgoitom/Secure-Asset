package entity

// User represents a user in the system.
type User struct {
	ID              string `bson:"_id,omitempty"`
	Username        string `bson:"username" validate:"required"`
	Email           string `bson:"email" validate:"required,email"`
	Password        string `bson:"password" validate:"required"`
	Role            string `bson:"role" validate:"required"` // For RBAC (e.g., "Admin", "Employee", "AssetManager")
	Department      string `bson:"department"`               // For ABAC
	ClearanceLevel  string `bson:"clearanceLevel"`           // For MAC (e.g., "Confidential", "Restricted", "Internal", "Public")
	IsMFAEnabled    bool   `bson:"isMFAEnabled"`
	OTPSecret       string `bson:"otpSecret,omitempty"`
	IsAccountLocked bool   `bson:"isAccountLocked"`
	// Token           string `bson:"token,omitempty"`
}
