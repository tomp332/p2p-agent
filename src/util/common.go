<<<<<<<< HEAD:src/util/common.go
package util
========
package common
>>>>>>>> origin/master:src/common/utils.go

import "github.com/google/uuid"

// GenerateRandomID random string generate
func GenerateRandomID() string {
	u, err := uuid.NewUUID()
	if err != nil {
		Logger.Fatal().Err(err).Msg("Could not generate UUID")
	}
	return u.String()
}
