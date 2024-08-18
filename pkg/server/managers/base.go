package managers

import "github.com/tomp332/p2p-agent/pkg/utils/configs"

type AuthenticationManager interface {
	Verify(accessToken string) (interface{}, error)
	Generate(username string, nodeType configs.NodeType) (string, error)
}
