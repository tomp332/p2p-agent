package test_p2p

import (
	"github.com/tomp332/p2p-agent/src/node/p2p"
	testUtils "github.com/tomp332/p2p-agent/tests/utils"
	"testing"
)

func TestInvalidNodeType(t *testing.T) {
	grpcTestServer := testUtils.NewTestAgentServer(t)
	config := testUtils.BaseNodeConfig(t)
	config.BaseNodeConfigs.Type = "invalid"
	p2p.NewBaseNode(grpcTestServer, &config.BaseNodeConfigs)
}
