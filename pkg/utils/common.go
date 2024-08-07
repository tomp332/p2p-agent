package utils

import (
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

// GenerateRandomID random string generate
func GenerateRandomID() string {
	u, err := uuid.NewUUID()
	if err != nil {
		log.Fatal().Err(err).Msg(err.Error())
	}
	return u.String()
}
