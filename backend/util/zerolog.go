package util

import (
	"context"
	"github.com/rs/zerolog"
)

func SetLogLevel(ctx context.Context, level zerolog.Level) context.Context {
	logger := zerolog.Ctx(ctx).Level(level)
	return logger.WithContext(ctx)
}
