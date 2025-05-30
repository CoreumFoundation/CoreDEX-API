package currency

import "github.com/CoreumFoundation/CoreDEX-API/utils/logger"

func (a *Application) schema() {
	_, err := a.client.Client.Exec(`CREATE TABLE IF NOT EXISTS Currency (
        Denom JSON,
        SendCommission JSON,
        BurnRate JSON,
        InitialAmount JSON,
        Chain VARCHAR(255),
        OriginChain VARCHAR(255),
        ChainSupply VARCHAR(255),
        Description VARCHAR(255),
        SkipDisplay BOOLEAN,
        MetaData JSON,
		Network INT,
		DenomString VARCHAR(255),
        UNIQUE KEY (DenomString, Network)
    )`)
	if err != nil {
		logger.Fatalf("Error creating Currency table: %v", err)
	}
}
