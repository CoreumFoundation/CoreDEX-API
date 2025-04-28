package trade

import "github.com/CoreumFoundation/CoreDEX-API/utils/logger"

// Initialize tables and indexes
func (a *Application) schema() {
	a.createTables()
	a.alterTables()
}

func (a *Application) createTables() {
	_, err := a.client.Client.Exec(`CREATE TABLE IF NOT EXISTS Trade (
		TXID VARCHAR(255),
		Account VARCHAR(255),
		OrderID VARCHAR(255),
		Sequence BIGINT,
		Amount JSON,
		Price FLOAT,
		Denom1 JSON,
		Denom2 JSON,
		Side INT,
		BlockTime JSON,
		BlockHeight BIGINT,
		MetaData JSON,
		USD FLOAT,
		Network INT,
		UNIQUE KEY (TXID, Sequence, Network)
	)`)
	if err != nil {
		logger.Fatalf("Error creating Trade table: %v", err)
	}
	_, err = a.client.Client.Exec(`CREATE TABLE IF NOT EXISTS TradePairs (
		Denom1 JSON,
		Denom2 JSON,
		MetaData JSON,
		Symbol1 VARCHAR(255) AS (JSON_UNQUOTE(JSON_EXTRACT(Denom1, '$.Denom'))),
		Symbol2 VARCHAR(255) AS (JSON_UNQUOTE(JSON_EXTRACT(Denom2, '$.Denom'))),
		Network INT AS (JSON_UNQUOTE(JSON_EXTRACT(MetaData, '$.Network'))),
		UNIQUE KEY (Symbol1, Symbol2, Network)
	)`)
	if err != nil {
		logger.Fatalf("Error creating TradePairs table: %v", err)
	}
}

func (a *Application) alterTables() {
	/*
		Create virtual columns for the Trade table for where indexes need to be created on:
		- Denom1.Denom
		- Denom2.Denom
		- BlockTime.seconds
	*/
	// Has to succeed or we first have to write logic to check if the columns exist....
	a.client.Client.Exec(`ALTER TABLE Trade 
	ADD COLUMN Symbol1 VARCHAR(255) AS (JSON_UNQUOTE(JSON_EXTRACT(Denom1, '$.Denom'))) STORED, 
	ADD COLUMN Symbol2 VARCHAR(255) AS (JSON_UNQUOTE(JSON_EXTRACT(Denom2, '$.Denom'))) STORED, 
	ADD COLUMN BlockTimeSeconds BIGINT AS (JSON_UNQUOTE(JSON_EXTRACT(BlockTime, '$.seconds'))) STORED`)
	a.client.Client.Exec(`ALTER TABLE TradePairs
	DROP COLUMN Symbol1,
	DROP COLUMN Symbol2,
	DROP COLUMN Network`)
	a.client.Client.Exec(`ALTER TABLE TradePairs 
	ADD COLUMN Currency1 VARCHAR(255) AS (JSON_UNQUOTE(JSON_EXTRACT(Denom1, '$.Currency'))) STORED, 
	ADD COLUMN Currency2 VARCHAR(255) AS (JSON_UNQUOTE(JSON_EXTRACT(Denom2, '$.Currency'))) STORED, 
	ADD COLUMN Issuer1 VARCHAR(255) AS (JSON_UNQUOTE(JSON_EXTRACT(Denom1, '$.Issuer'))) STORED, 
	ADD COLUMN Issuer2 VARCHAR(255) AS (JSON_UNQUOTE(JSON_EXTRACT(Denom2, '$.Issuer'))) STORED, 
	ADD COLUMN Network INT AS (JSON_UNQUOTE(JSON_EXTRACT(MetaData, '$.Network'))) STORED`)
	// Add unique key to TradePairs table:
	a.client.Client.Exec(`ALTER TABLE TradePairs
	ADD UNIQUE KEY (Currency1, Currency2, Issuer1, Issuer2, Network)`)
	// Addition of enriched field to have a flexible skip of temporary failures:
	a.client.Client.Exec(`ALTER TABLE Trade
	ADD COLUMN Enriched BOOLEAN DEFAULT TRUE`)
	a.client.Client.Exec(`ALTER TABLE Trade
	ADD COLUMN Inverted BOOLEAN DEFAULT FALSE`)
	a.client.Client.Exec(`ALTER TABLE TradePairs 
	ADD Column PriceTick BIGINT,
	ADD COLUMN QuantityStep INT`)
	a.client.Client.Exec(`ALTER TABLE TradePairs 
	DROP Column PriceTick`)
	a.client.Client.Exec(`ALTER TABLE TradePairs 
	ADD Column PriceTick JSON`)
}

func (a *Application) index() {
	a.client.Client.Exec(`DROP INDEX trade_1 ON Trade`)
	a.client.Client.Exec(`CREATE INDEX trade_3 ON Trade (
		Symbol1,
		Symbol2,
		BlockTimeSeconds,
		Network,
		Side
	)`)
	a.client.Client.Exec(`CREATE INDEX trade_2 ON Trade (
		Account,
		Symbol1,
		Symbol2,
		BlockTimeSeconds,
		Network
	)`)

	a.client.Client.Exec(`CREATE INDEX tradepairs_1 ON TradePairs (
		Currency1(50),
		Currency2(50),
		Issuer1(50),
		Issuer2(50),
		Network
	)`)
	a.client.Client.Exec(`CREATE INDEX tradepairs_2 ON TradePairs (
		Network
	)`)
}
