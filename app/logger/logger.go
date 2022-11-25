package logger

import (
	"io"
	"os"
	"time"

	"github.com/ilhamfzri/pendek.in/config"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
	"gopkg.in/natefinch/lumberjack.v2"
)

type Logger struct {
	*zerolog.Logger
}

func NewLogger(cfg config.LoggerConfig) *Logger {
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	zerolog.TimeFieldFormat = time.RFC3339Nano

	var output io.Writer = zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: time.RFC3339,
	}

	if os.Getenv("APP_STAGE") != "development" {
		fileLogger := &lumberjack.Logger{
			Filename:   cfg.Output,
			MaxSize:    100, //
			MaxBackups: 10,
			MaxAge:     14,
			Compress:   true,
		}
		output = zerolog.MultiLevelWriter(os.Stderr, fileLogger)
	}

	logger := zerolog.New(output).
		Level(zerolog.Level(cfg.Level)).
		With().
		Timestamp().
		Logger()

	return &Logger{&logger}
}
func (l *Logger) FatalIfErr(err error, message string) {
	if err != nil {
		l.Fatal().Err(err).Msg(message)
	}
}

func (l *Logger) PanicIfErr(err error, message string) {
	if err != nil {
		l.Panic().Err(err).Msg(message)
	}
}
