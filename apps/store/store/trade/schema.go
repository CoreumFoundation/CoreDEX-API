package trade

import (
	"database/sql"

	"github.com/CoreumFoundation/CoreDEX-API/utils/logger"
)

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
		Denom1 JSON DEFAULT NULL,
		Denom2 JSON DEFAULT NULL,
		MetaData JSON DEFAULT NULL,
		Currency1  VARCHAR(100) AS (JSON_UNQUOTE(JSON_EXTRACT(Denom1, '$.Currency'))) STORED, 
		Currency2 VARCHAR(100) AS (JSON_UNQUOTE(JSON_EXTRACT(Denom2, '$.Currency'))) STORED, 
		Issuer1 VARCHAR(100) AS (JSON_UNQUOTE(JSON_EXTRACT(Denom1, '$.Issuer'))) STORED, 
		Issuer2 VARCHAR(100) AS (JSON_UNQUOTE(JSON_EXTRACT(Denom2, '$.Issuer'))) STORED,
		Network  INT AS (JSON_UNQUOTE(JSON_EXTRACT(MetaData, '$.Network'))) STORED,
		QuantityStep INT DEFAULT NULL,
		PriceTick JSON DEFAULT NULL,
		UNIQUE KEY tradepairs_1 (Currency1,Currency2,Issuer1,Issuer2,Network),
		KEY tradepairs_2 (Network))`)
	if err != nil {
		logger.Fatalf("Error creating TradePairs table: %v", err)
	}

	// Check if we need to perform the conversion by looking for NULL values
	var nullCount int
	err = a.client.Client.QueryRow(`SELECT COUNT(*) FROM (SELECT COUNT(*) FROM TradePairs 
	WHERE PriceTick IS NOT NULL 
	GROUP BY Denom1, Denom2, QuantityStep, PriceTick
	HAVING COUNT(*) > 1) AS t`).Scan(&nullCount)
	if err != nil && err != sql.ErrNoRows {
		logger.Fatalf("Error checking for NULL values: %v", err)
	}
	if nullCount > 0 {
		// Fix for the NULL values being parsed into the Currency1, Currency2, Issuer1, Issuer2 columns:
		_, err = a.client.Client.Exec(`CREATE TABLE TradePairs_duplicate_fix LIKE TradePairs;
			INSERT INTO TradePairs_duplicate_fix (Denom1, Denom2, MetaData, QuantityStep, PriceTick)
			SELECT DISTINCT Denom1, Denom2, MetaData, QuantityStep, PriceTick 
			FROM TradePairs;

			ALTER TABLE TradePairs_duplicate_fix
				MODIFY COLUMN Currency1 VARCHAR(100) AS (COALESCE(JSON_UNQUOTE(JSON_EXTRACT(Denom1, '$.Currency')), '')) STORED,
				MODIFY COLUMN Currency2 VARCHAR(100) AS (COALESCE(JSON_UNQUOTE(JSON_EXTRACT(Denom2, '$.Currency')), '')) STORED,
				MODIFY COLUMN Issuer1 VARCHAR(100) AS (COALESCE(JSON_UNQUOTE(JSON_EXTRACT(Denom1, '$.Issuer')), '')) STORED,
				MODIFY COLUMN Issuer2 VARCHAR(100) AS (COALESCE(JSON_UNQUOTE(JSON_EXTRACT(Denom2, '$.Issuer')), '')) STORED,
				ADD UNIQUE KEY tradepairs_1 (Currency1,Currency2,Issuer1,Issuer2,Network),
				ADD KEY tradepairs_2 (Network);

			DROP TABLE TradePairs;
			
			RENAME TABLE TradePairs_duplicate_fix TO TradePairs`)
		if err != nil {
			logger.Errorf("Error fixing TradePairs table: %v", err)
			// Reinit the table but now with the correct columns:
			_, err = a.client.Client.Exec(`DROP TABLE TradePairs`)
			if err != nil {
				logger.Errorf("Error dropping TradePairs table: %v", err)
			}
			_, err = a.client.Client.Exec(`CREATE TABLE IF NOT EXISTS TradePairs (
				Denom1 JSON DEFAULT NULL,
				Denom2 JSON DEFAULT NULL,
				MetaData JSON DEFAULT NULL,
				Currency1  VARCHAR(100) AS (COALESCE(JSON_UNQUOTE(JSON_EXTRACT(Denom1, '$.Currency')), '')) STORED,
				Currency2 VARCHAR(100) AS (COALESCE(JSON_UNQUOTE(JSON_EXTRACT(Denom2, '$.Currency')), '')) STORED,
				Issuer1 VARCHAR(100) AS (COALESCE(JSON_UNQUOTE(JSON_EXTRACT(Denom1, '$.Issuer')), '')) STORED,
				Issuer2 VARCHAR(100) AS (COALESCE(JSON_UNQUOTE(JSON_EXTRACT(Denom2, '$.Issuer')), '')) STORED,
				Network  INT AS (JSON_UNQUOTE(JSON_EXTRACT(MetaData, '$.Network'))) STORED,
				QuantityStep INT DEFAULT NULL,
				PriceTick JSON DEFAULT NULL,
				UNIQUE KEY tradepairs_1 (Currency1,Currency2,Issuer1,Issuer2,Network),
				KEY tradepairs_2 (Network))`)
			if err != nil {
				logger.Fatalf("Error creating TradePairs table: %v", err)
			}
		}
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
	// Addition of enriched field to have a flexible skip of temporary failures:
	a.client.Client.Exec(`ALTER TABLE Trade
	ADD COLUMN Enriched BOOLEAN DEFAULT TRUE`)
	a.client.Client.Exec(`ALTER TABLE Trade
	ADD COLUMN Inverted BOOLEAN DEFAULT FALSE`)
	// Add Processed column to Trade table:
	a.client.Client.Exec(`ALTER TABLE Trade
	ADD COLUMN Processed BOOLEAN DEFAULT FALSE`)
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
	a.client.Client.Exec(`CREATE INDEX tradepairs_3 ON TradePairs (
		Network,
		Currency1,
		Currency2,
		Issuer1,
		Issuer2
	)`)
}
