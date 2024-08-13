package test_p2p

import (
	"github.com/stretchr/testify/assert"
	"github.com/tomp332/p2p-agent/src/node"
	"github.com/tomp332/p2p-agent/src/node/p2p"
	"github.com/tomp332/p2p-agent/src/utils/configs"
	testUtils "github.com/tomp332/p2p-agent/tests/utils"
	"testing"
)

func TestInvalidNodeType(t *testing.T) {
	_ = testUtils.NewTestAgentServer(t)
	config := configs.NodeConfigs{}
	config.Type = "invalid"
	p2p.NewBaseNode(&config)
}

func TestDuplicateNodeTypes(t *testing.T) {
	server := testUtils.NewTestAgentServer(t)
	config1 := configs.NodeConfigs{}
	config1.Type = "file"
	config2 := configs.NodeConfigs{}
	config2.Type = "file"
	nodeConfigs := []configs.NodeConfigs{config1, config2}
	_, err := node.InitializeNodes(server, &nodeConfigs)
	assert.NotEmpty(t, err)
}
