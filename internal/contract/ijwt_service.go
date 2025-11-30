package contract

import (
	"github.com/mikiasgoitom/Secure-Asset/internal/domain/entity"
	"github.com/mikiasgoitom/Secure-Asset/internal/dto"
)

type IJWTService interface {
	GenerateToken(user *entity.User) (string, error)
	ValidateToken(tokenString string) (*dto.CustomClaims, error)
}