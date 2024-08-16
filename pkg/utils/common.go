package utils

import (
	"crypto/rand"
	"encoding/base64"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
)

// GenerateRandomID random string generate
func GenerateRandomID() string {
	u, err := uuid.NewUUID()
	if err != nil {
		log.Fatal().Err(err).Msg(err.Error())
	}
	return u.String()
}

// GenerateSecureToken generates a secure random token of the specified byte length
func GenerateSecureToken(length int) string {
	// Create a byte slice of the specified length
	token := make([]byte, length)

	// Fill the slice with secure random bytes
	_, err := rand.Read(token)
	if err != nil {
		log.Fatal().Err(err).Msg(err.Error())
	}

	// Encode the byte slice into a base64 string to make it URL-safe
	return base64.URLEncoding.EncodeToString(token)
}

func BcryptCompare(hashedPassword string, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}
