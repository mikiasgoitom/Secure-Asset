package logger

import (
	"github.com/mikiasgoitom/Secure-Asset/internal/contract"
	"github.com/mikiasgoitom/Secure-Asset/internal/domain/valueobject"
	"go.uber.org/zap"
)

type ZapAdapter struct{
	logger *zap.Logger
}

func NewZapAdapter(isProduction bool) (contract.ILogger, error) {
	var logger *zap.Logger
	var err error
	if isProduction {
		logger, err = zap.NewProduction()
	} else {
		logger, err = zap.NewDevelopment()
	}
	if err != nil {
		return nil, err
	}
	return &ZapAdapter{logger: logger}, nil
}

func (z *ZapAdapter) toZapFields(fields ...valueobject.Field) []zap.Field {
	zapFields := make([]zap.Field, len(fields))
	for i, field := range fields {
		zapFields[i] = zap.Any(field.Key, field.Value)
	}
	return zapFields
}

func (z *ZapAdapter) Info(msg string, fields ...valueobject.Field) {
	z.logger.Info(msg, z.toZapFields(fields...)...)
}
func (z *ZapAdapter) Error(msg string, fields ...valueobject.Field) {
	z.logger.Error(msg, z.toZapFields(fields...)...)
}
func (z *ZapAdapter) Debug(msg string, fields ...valueobject.Field) {
	z.logger.Debug(msg, z.toZapFields(fields...)...)
}
func (z *ZapAdapter) Warn(msg string, fields ...valueobject.Field) {
	z.logger.Warn(msg, z.toZapFields(fields...)...)
}
func (z *ZapAdapter) Fatal(msg string, fields ...valueobject.Field) {
	z.logger.Fatal(msg, z.toZapFields(fields...)...)
}

func (z *ZapAdapter) Sync() error {
	return z.logger.Sync()
}