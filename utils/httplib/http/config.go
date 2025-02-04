package http

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/CoreumFoundation/CoreDEX-API/utils/logger"
)

const EnvConfig = "HTTP_CONFIG"

type Duration struct {
	time.Duration
}

type httpConfig struct {
	CORS    *CorsConfig `json:"cors"`
	Port    string      `json:"port"`
	TimeOut *TimeOut    `json:"timeouts"`
}

type CorsConfig struct {
	AllowedOrigins []string `json:"allowedOrigins"`
}

type TimeOut struct {
	Read     Duration `json:"read"`
	Write    Duration `json:"write"`
	Idle     Duration `json:"idle"`
	Shutdown Duration `json:"shutdown"`
}

/*
parseConfig parses the configuration yaml and sets the configuration
If the parsing fails, the application will exit with a fatal
Config is expected to be present as environment variable "HTTP_CONFIG" : see README.md!
*/
func parseConfig() {
	cfg := os.Getenv(EnvConfig)
	logger.Infof("Found %s: %s", EnvConfig, cfg)
	if cfg == "" {
		logger.Fatalf("%s env is not set", EnvConfig)
	}
	// Parse the config:
	if err := json.Unmarshal([]byte(cfg), &conf); err != nil {
		logger.Fatalf("Could not parse %s: %v", EnvConfig, err)
	}
}

func (duration *Duration) UnmarshalJSON(b []byte) error {
	var unmarshalledJson interface{}

	err := json.Unmarshal(b, &unmarshalledJson)
	if err != nil {
		return err
	}

	switch value := unmarshalledJson.(type) {
	case float64:
		duration.Duration = time.Duration(value)
	case string:
		duration.Duration, err = time.ParseDuration(value)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("invalid duration: %#v", unmarshalledJson)
	}

	return nil
}
