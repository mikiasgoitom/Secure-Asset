package contract

import "github.com/mikiasgoitom/Secure-Asset/internal/domain/valueobject"

type ILogger interface {
	Info(msg string, fields ...valueobject.Field) 
	Error(msg string, fields ...valueobject.Field)
	Debug(msg string, fields ...valueobject.Field)
	Fatal(msg string, fields ...valueobject.Field)
	Warn(msg string, fields ...valueobject.Field)
	Sync() error
}