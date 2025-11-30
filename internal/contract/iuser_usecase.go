package contract

import (
	"context"

	"github.com/mikiasgoitom/Secure-Asset/internal/domain/entity"
)

type IUserUsecase interface {
	Register(ctx context.Context, username, email, password string) (*entity.User, error)
	Login(ctx context.Context, identifier, password string) (string, error)
}
