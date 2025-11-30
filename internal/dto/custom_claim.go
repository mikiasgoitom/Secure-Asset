package dto

import "github.com/golang-jwt/jwt/v5"

type CustomClaims struct {
	UserID         string `json:"userID"`
	Role           string `json:"role"`
	Department     string `json:"department"`
	ClearanceLevel string `json:"clearanceLevel"`
	jwt.RegisteredClaims
}
