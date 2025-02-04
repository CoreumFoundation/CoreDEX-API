package domain

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/CoreumFoundation/CoreDEX-API/domain/metadata"
)

func GetTokenRegistryEntries(ctx context.Context, network metadata.Network) (map[string]Asset, error) {
	assets := make(map[string]Asset)
	jsonFile := fmt.Sprintf(
		"https://github.com/CoreumFoundation/token-registry/raw/refs/heads/master/%s/assets.json",
		strings.ToLower(network.String()),
	)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, jsonFile, http.NoBody)
	if err != nil {
		return assets, err
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return assets, err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return assets, fmt.Errorf("status code %d", res.StatusCode)
	}
	var tokenRegistryAssets TokenRegistryAssets
	if err = json.NewDecoder(res.Body).Decode(&tokenRegistryAssets); err != nil {
		return assets, err
	}
	for _, asset := range tokenRegistryAssets.Assets {
		assets[asset.Denom] = asset
	}
	return assets, nil
}

type IbcInfo struct {
	DisplayName string `json:"display_name,omitempty"`
	Precision   int    `json:"precision,omitempty"`
	SourceChain string `json:"source_chain,omitempty"`
	Denom       string `json:"denom"`
}

type LogoURIs struct {
	Png string `json:"png"`
	Svg string `json:"svg"`
}

type XrplInfo struct {
	Precision   int    `json:"precision"`
	SourceChain string `json:"source_chain"`
	Issuer      string `json:"issuer"`
	Currency    string `json:"currency"`
}

type Extra struct {
	XrplInfo XrplInfo `json:"xrpl_info,omitempty"`
	IbcInfo  IbcInfo  `json:"ibc_info,omitempty"`
}

type Asset struct {
	Denom       string   `json:"denom"`
	TokenName   string   `json:"token_name"`
	Decimals    int      `json:"decimals"`
	Description string   `json:"description"`
	IbcInfo     IbcInfo  `json:"ibc_info"`
	LogoURIs    LogoURIs `json:"logo_URIs"`
	Extra       Extra    `json:"extra"`
}

type TokenRegistryAssets struct {
	Assets []Asset `json:"assets"`
}
