package coreum

import (
	"encoding/json"
	"os"

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
	return v
}
