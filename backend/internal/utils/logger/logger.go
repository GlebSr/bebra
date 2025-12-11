package logger

import (
	"context"
	"os"

	"github.com/rs/zerolog"
)

type coreLogger struct {
	zl zerolog.Logger
}

var log coreLogger

func InitLogger(env string) {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.SetGlobalLevel(zerolog.TraceLevel)

	var output zerolog.Logger
	if env == "prod" {
		output = zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout})
	} else {
		output = zerolog.New(os.Stdout)
	}

	log.zl = output.With().
		Timestamp().
		Logger()
}

func (l *coreLogger) logEvent(ctx context.Context, event *zerolog.Event, format string, args ...interface{}) {
	if ctx != nil {
		params, ok := ctx.Value("fields_to_log").(map[string]any)
		if ok {
			event.Fields(params)
		}
	}
	event.Msgf(format, args...)
}

func Debugf(ctx context.Context, format string, args ...interface{}) {
	log.logEvent(ctx, log.zl.Debug(), format, args...)
}

func Infof(ctx context.Context, format string, args ...interface{}) {
	log.logEvent(ctx, log.zl.Info(), format, args...)
}

func Warnf(ctx context.Context, format string, args ...interface{}) {
	log.logEvent(ctx, log.zl.Warn(), format, args...)
}

func Errorf(ctx context.Context, format string, args ...interface{}) {
	log.logEvent(ctx, log.zl.Error(), format, args...)
}

func Fatalf(ctx context.Context, format string, args ...interface{}) {
	log.logEvent(ctx, log.zl.Fatal(), format, args...)
}
