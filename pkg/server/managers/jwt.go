package managers

import (
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/rs/zerolog/log"
	"github.com/tomp332/p2p-agent/pkg/utils/configs"
	"strings"
	"time"
)

type JWTManager struct {
	secretKey     string
	tokenDuration time.Duration
}

type PeerToken struct {
	jwt.StandardClaims
	Username string           `json:"username"`
	NodeType configs.NodeType `json:"node_type"`
}

func NewJWTManager(secretKey string, tokenDuration time.Duration) *JWTManager {
	return &JWTManager{secretKey, tokenDuration}
}

func (manager *JWTManager) Generate(username string, nodeType configs.NodeType) (string, error) {
	claims := PeerToken{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(manager.tokenDuration).Unix(),
		},
		Username: username,
		NodeType: nodeType,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(manager.secretKey))
}

func (manager *JWTManager) Verify(accessToken string) (interface{}, error) {
	tokenOnly, err := extractBearerToken(accessToken)
	if err != nil {
		log.Error().Err(err).Msg("Error extracting token")
		return nil, fmt.Errorf("invalid token: %w", err)
	}
	token, err := jwt.ParseWithClaims(
		tokenOnly,
		&PeerToken{},
		func(token *jwt.Token) (interface{}, error) {
			_, ok := token.Method.(*jwt.SigningMethodHMAC)
			if !ok {
				return nil, fmt.Errorf("unexpected token signing method")
			}
			return []byte(manager.secretKey), nil
		},
	)

	if err != nil {
		log.Error().Err(err).Msg("Error parsing jwt token.")
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	claims, ok := token.Claims.(*PeerToken)
	if !ok {
		log.Error().Err(err).Msg("Error validating auth token")
		return nil, fmt.Errorf("invalid token claims")
	}

	return claims, nil
}

func extractBearerToken(accessToken string) (string, error) {
	// Check if the accessToken starts with the expected "Bearer " prefix
	const bearerPrefix = "Bearer "

	if !strings.HasPrefix(accessToken, bearerPrefix) {
		return "", errors.New("invalid authorization token format")
	}

	// Extract the token after the "Bearer " prefix
	token := strings.TrimSpace(strings.TrimPrefix(accessToken, bearerPrefix))

	// Check if the extracted token is not empty
	if token == "" {
		return "", errors.New("authorization token is empty")
	}

	return token, nil
}
