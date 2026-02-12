package logger

import (
	"bytes"
	"encoding/json"
	"os"

	"github.com/rs/zerolog"
)

type prettyJSONWriter struct{}

func (w prettyJSONWriter) Write(p []byte) (n int, err error) {
	var m map[string]interface{}
	err = json.Unmarshal(p, &m)
	if err != nil {
		return 0, err
	}

	level := m["level"]

	var colorStart string
	colorReset := "\033[0m"

	switch level {
	case "error":
		colorStart = "\033[31m" // red
	case "warn":
		colorStart = "\033[33m" // yellow
	case "info":
		colorStart = "\033[34m" // blue
	case "debug":
		colorStart = "\033[32m" // green
	default:
		colorStart = ""
	}

	var out bytes.Buffer
	err = json.Indent(&out, p, "", "  ")
	if err != nil {
		return 0, err
	}

	return os.Stdout.Write([]byte(colorStart + out.String() + colorReset))
}

func New(env string) zerolog.Logger {
	if env == "PRODUCTION" {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)

		return zerolog.New(os.Stdout).
			With().
			Timestamp().
			Logger()
	}

	zerolog.SetGlobalLevel(zerolog.DebugLevel)

	return zerolog.New(prettyJSONWriter{}).
		With().
		Timestamp().
		Logger()
}
