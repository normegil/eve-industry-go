package config

import (
	"fmt"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
	"strconv"
)

func init() {
	log.Logger = zerolog.New(os.Stdout).With().Timestamp().Logger()
	colorizedLogStr := os.Getenv(appPrefix + "COLORIZED_LOG")
	colorizedLog := false
	if len(colorizedLogStr) != 0 {
		var err error
		colorizedLog, err = strconv.ParseBool(colorizedLogStr)
		if err != nil {
			panic(fmt.Errorf("invalid bool '%s': %w", colorizedLogStr, err))
		}
	}
	if colorizedLog {
		log.Logger = log.Logger.Output(zerolog.ConsoleWriter{Out: os.Stdout})
	}
}
