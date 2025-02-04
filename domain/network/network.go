package network

import (
	"errors"
	"net/http"
	"strings"

	"github.com/CoreumFoundation/CoreDEX-API/domain/metadata"
)

const networkKey = "Network"

var ErrInvalidNetwork = errors.New("invalid network")

// Helper function to retrieve the network from the request header in a consistent fashion
func Network(r *http.Request) (metadata.Network, error) {
	// Get the network from the headers.
	// Check against the enum auth.Network (case insensitive).
	network := r.Header.Get(networkKey)
	// check against the enum auth.Network
	if _, ok := metadata.Network_value[strings.ToUpper(network)]; !ok {
		return metadata.Network_NETWORK_DO_NOT_USE, ErrInvalidNetwork
	}
	return metadata.Network(metadata.Network_value[strings.ToUpper(network)]), nil
}
