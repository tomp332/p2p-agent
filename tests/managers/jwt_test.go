package managers

import (
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/assert"
	"github.com/tomp332/p2p-agent/pkg/server/managers"
	"github.com/tomp332/p2p-agent/pkg/utils/configs"
)

func Test_Generate(t *testing.T) {
	secretKey := "test-secret-key"
	tokenDuration := time.Minute * 5
	jwtManager := managers.NewJWTManager(secretKey, tokenDuration)

	t.Run("Generate Token Successfully", func(t *testing.T) {
		username := "validUser"
		nodeType := configs.FilesNodeType

		token, err := jwtManager.Generate(username, nodeType)
		assert.NoError(t, err)
		assert.NotEmpty(t, token)
	})

	t.Run("Generate Token and Verify", func(t *testing.T) {
		username := "validUser"
		nodeType := configs.FilesNodeType

		token, err := jwtManager.Generate(username, nodeType)
		assert.NoError(t, err)

		// Verify the token
		claims, err := jwtManager.Verify("Bearer " + token)
		assert.NoError(t, err)

		peerToken, ok := claims.(*managers.PeerToken)
		assert.True(t, ok)
		assert.Equal(t, username, peerToken.Username)
		assert.Equal(t, nodeType, peerToken.NodeType)
	})
}

func Test_Verify(t *testing.T) {
	secretKey := "test-secret-key"
	tokenDuration := time.Minute * 5
	jwtManager := managers.NewJWTManager(secretKey, tokenDuration)

	t.Run("Invalid Token Format", func(t *testing.T) {
		_, err := jwtManager.Verify("InvalidToken")
		assert.Error(t, err)
		assert.Equal(t, "invalid token: invalid authorization token format", err.Error())
	})

	t.Run("Empty Token", func(t *testing.T) {
		_, err := jwtManager.Verify("Bearer ")
		assert.Error(t, err)
		assert.Equal(t, "invalid token: authorization token is empty", err.Error())
	})

	t.Run("Invalid Signing Method", func(t *testing.T) {
		invalidToken := jwt.NewWithClaims(jwt.SigningMethodRS256, &managers.PeerToken{
			Username: "validUser",
			NodeType: configs.FilesNodeType,
		})
		tokenString, _ := invalidToken.SignedString([]byte(secretKey))

		_, err := jwtManager.Verify("Bearer " + tokenString)
		assert.Error(t, err)
		assert.Equal(t, "invalid token: authorization token is empty", err.Error())
	})
}
