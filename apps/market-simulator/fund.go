package main

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/CoreumFoundation/CoreDEX-API/utils/logger"
)

/*
Fund uses the address of the user to fund the account with the specified amount udevcore

The call takes a post:

https://api.devnet-1.coreum.dev/api/faucet/v1/fund
{address: "devcore19p7572k4pj00szx36ehpnhs8z2gqls8ky3ne43"}
*/

type fund struct {
	Address string `json:"address"`
}

func addFunds(address string) {
	url := "https://api.devnet-1.coreum.dev/api/faucet/v1/fund"
	payload := fund{
		Address: address,
	}
	jsonData, err := json.Marshal(payload)
	if err != nil {
		logger.Errorf("Error marshalling JSON:", err)
		return
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		logger.Errorf("Error creating request:", err)
		return
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logger.Errorf("Error making request:", err)
		return
	}
	defer resp.Body.Close()
	logger.Infof("Response status: %s", resp.Status)
}
