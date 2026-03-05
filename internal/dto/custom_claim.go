package dto

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/mikiasgoitom/Secure-Asset/internal/domain/valueobject"
)

type CustomClaims struct {
	UserID         string                     `json:"userID"`
	Role           string                     `json:"role"`
	Department     string                     `json:"department"`
	ClearanceLevel valueobject.Classification `json:"clearanceLevel"`
	jwt.RegisteredClaims
}
