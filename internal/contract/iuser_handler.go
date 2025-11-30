package contract

import (
	"context"

	"github.com/mikiasgoitom/Secure-Asset/internal/domain/entity"
)

type IUserHandler interface {
	RegisterUser(ctx context.Context) (*entity.User, error)
}