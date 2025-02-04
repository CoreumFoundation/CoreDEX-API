package rates

import (
	"fmt"
)

// local value - coingecko value
var mapping = map[string]string{
	"UCORE": "coreum",
}

var reverseMapping = reverseMap(mapping)

func assetIDToCoingeckoID(assetID string) (string, error) {
	if coingeckoID, ok := mapping[normalizeAssetID(assetID)]; ok {
		return coingeckoID, nil
	}
	return "", fmt.Errorf("unsupported assetID provided: %v", assetID)
}

func coingeckoIDToAssetID(coingeckoID string) (string, error) {
	if assetID, ok := reverseMapping[coingeckoID]; ok {
		return assetID, nil
	}
	return "", fmt.Errorf("unsupported coingeckoID provided: %v", coingeckoID)
}

// TODO: Multi currency in this function does not work when it has the same assetID
func reverseMap(m map[string]string) map[string]string {
	n := make(map[string]string, len(m))
	for k, v := range m {
		if _, ok := n[v]; !ok {
			n[v] = k
		}
	}
	return n
}
