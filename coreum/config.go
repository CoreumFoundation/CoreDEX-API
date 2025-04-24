package coreum

import (
	"encoding/json"
	"os"
	"strings"

	"github.com/CoreumFoundation/CoreDEX-API/domain/metadata"
	"github.com/CoreumFoundation/CoreDEX-API/utils/logger"
)

type Config struct {
	Node []Node `json:"Node"`
}

type Node struct {
	Network  string
	GRPCHost string
	RPCHost  string
}

func ParseConfig() *Config {
	conf := os.Getenv("NETWORKS")
	if conf == "" {
		logger.Fatalf("NETWORKS environment variable not set")
	}
	v := &Config{}
	err := json.Unmarshal([]byte(conf), v)
	if err != nil {
		logger.Fatalf("failed to parse NETWORKS config: %v", err)
	}
	// Validate networks in config:
	// - Network name should be one of the known networks in the enum metadata.Network
	for _, node := range v.Node {
		if metadata.Network_value[strings.ToUpper(node.Network)] == 0 {
			logger.Fatalf("invalid network name: %s", node.Network)
		}
	}
	return v
}
