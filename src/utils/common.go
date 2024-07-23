package utils

import (
	"encoding/json"
	"github.com/google/uuid"
)

// GenerateRandomID random string generate
func GenerateRandomID() string {
	u, err := uuid.NewUUID()
	if err != nil {
		HandleCriticalError(err)
	}
	return u.String()
}

// MapToStruct converts a map to a struct using JSON serialization
func MapToStruct[T any](data map[string]interface{}) (*T, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	var result T
	if err := json.Unmarshal(jsonData, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func HandleCriticalError(err error) {
	if err != nil {
		Logger.Fatal().Err(err).Msg(err.Error())
	}
}
