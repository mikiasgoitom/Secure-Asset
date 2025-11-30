package security

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/mikiasgoitom/Secure-Asset/internal/domain/entity"
	"github.com/mikiasgoitom/Secure-Asset/internal/dto"
)

type JWTService struct {
	sercretKey string
	issuer     string
}

func NewJWTService(secretKey, issuer string) *JWTService {
	return &JWTService{
		sercretKey: secretKey,
		issuer:     issuer,
	}
}

func (j *JWTService) GenerateToken(user *entity.User) (string, error) {
	claims := dto.CustomClaims{
		UserID:         user.ID,
		Role:           user.Role,
		Department:     user.Department,
		ClearanceLevel: user.ClearanceLevel,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24).UTC()),
			IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
			Issuer:    j.issuer,
			Subject:   user.ID,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.sercretKey))
}

func (j *JWTService) ValidateToken(tokenString string) (*dto.CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &dto.CustomClaims{}, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected sigining method: %v", t.Header["alg"])
		}
		return []byte(j.sercretKey), nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	if claims, ok := token.Claims.(*dto.CustomClaims); ok && token.Valid {
        return claims, nil
    }

	return nil, fmt.Errorf("invalid token claims")
}