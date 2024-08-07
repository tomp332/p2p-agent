package utils

import (
	"github.com/google/uuid"
)

// GenerateRandomID random string generate
func GenerateRandomID() string {
	u, err := uuid.NewUUID()
	if err != nil {
		Logger.Fatal().Err(err).Msg(err.Error())
	}
	return u.String()
}
