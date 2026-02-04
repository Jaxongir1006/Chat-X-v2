package logger

import (
	"os"

	"github.com/rs/zerolog"
)

func New(env string) zerolog.Logger {
	if env == "prod" {
		return zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout}).With().Timestamp().Logger()
	}

	return zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout}).With().Timestamp().Logger()
}
