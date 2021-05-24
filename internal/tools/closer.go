package tools

import (
	"github.com/rs/zerolog/log"
	"io"
)

func Close(closer io.Closer) {
	if err := closer.Close(); nil != err {
		log.Error().Err(err).Msg("Could not close resource")
	}
}
