package order

import "github.com/CoreumFoundation/CoreDEX-API/utils/logger"

func (a *Application) schema() {
	a.createTables()
	a.alterTables()
}

// Initialize tables and indexes
func (a *Application) createTables() {
	_, err := a.client.Client.Exec(`CREATE TABLE IF NOT EXISTS OrderData (
		Account VARCHAR(255),
		Type INT,
		OrderID VARCHAR(255),
		Sequence BIGINT,
		BaseDenom JSON,
		QuoteDenom JSON,
		Price DOUBLE,
		Quantity JSON,
		RemainingQuantity JSON,
		Side INT,
		GoodTil JSON,
		TimeInForce INT,
		BlockTime JSON,
		OrderFee BIGINT,
		MetaData JSON,
		TXID VARCHAR(255),
		BlockHeight BIGINT,
		Network INT,
		UNIQUE KEY (Sequence, Network)
	)`)
	if err != nil {
		logger.Fatalf("Error creating state table: %v", err)
	}
	// Create historical table and trigger
	_, err = a.client.Client.Exec(`
	CREATE TABLE IF NOT EXISTS OrderDataHistory (
		Account VARCHAR(255),
		Type INT,
		OrderID VARCHAR(255),
		Sequence BIGINT,
		BaseDenom JSON,
		QuoteDenom JSON,
		Price DOUBLE,
		Quantity JSON,
		RemainingQuantity JSON,
		Side INT,
		GoodTil JSON,
		TimeInForce INT,
		BlockTime JSON,
		OrderFee BIGINT,
		MetaData JSON,
		TXID VARCHAR(255),
		BlockHeight BIGINT,
		Network INT
	)`)
	if err != nil {
		logger.Fatalf("Error creating historical table OrderDataHistory: %v", err)
	}
}

func (a *Application) alterTables() {
	a.client.Client.Exec(`ALTER TABLE OrderData 
	ADD COLUMN BaseCurrency VARCHAR(255) AS (JSON_UNQUOTE(JSON_EXTRACT(BaseDenom, '$.Currency'))) STORED, 
	ADD COLUMN QuoteCurrency VARCHAR(255) AS (JSON_UNQUOTE(JSON_EXTRACT(QuoteDenom, '$.Currency'))) STORED, 
	ADD COLUMN BaseIssuer VARCHAR(255) AS (JSON_UNQUOTE(JSON_EXTRACT(BaseDenom, '$.Issuer'))) STORED, 
	ADD COLUMN QuoteIssuer VARCHAR(255) AS (JSON_UNQUOTE(JSON_EXTRACT(QuoteDenom, '$.Issuer'))) STORED, 
	ADD COLUMN BlockTimeSeconds BIGINT AS (JSON_UNQUOTE(JSON_EXTRACT(BlockTime, '$.seconds'))) STORED`)

	// Add the OrderStatus INT column (ignore error if it already exists)
	a.client.Client.Exec(`ALTER TABLE OrderData ADD COLUMN OrderStatus INT`)
	a.client.Client.Exec(`ALTER TABLE OrderDataHistory ADD COLUMN OrderStatus INT`)
	// Replace the trigger with the new one
	_, err := a.client.Client.Exec(`DROP TRIGGER IF EXISTS after_order_update`)
	if err != nil {
		logger.Fatalf("Error dropping trigger after_order_update: %v", err)
	}
	_, err = a.client.Client.Exec(`
	CREATE TRIGGER IF NOT EXISTS after_order_update
	AFTER UPDATE ON OrderData
	FOR EACH ROW
	BEGIN
		INSERT INTO OrderDataHistory (` + OrderDataFields + `) VALUES (
			NEW.Account,
			NEW.Type,
			NEW.OrderID,
			NEW.Sequence,
			NEW.BaseDenom,
			NEW.QuoteDenom,
			NEW.Price,
			NEW.Quantity,
			NEW.RemainingQuantity,
			NEW.Side,
			NEW.GoodTil,
			NEW.TimeInForce,
			NEW.BlockTime,
			NEW.OrderFee,
			NEW.MetaData,
			NEW.TXID,
			NEW.BlockHeight,
			NEW.OrderStatus,
			NEW.Network
		);
	END;`)
	if err != nil {
		logger.Fatalf("Error creating historical table OrderDataHistory: %v", err)
	}
}

func (a *Application) index() {
	a.client.Client.Exec(`CREATE INDEX orderdata_1 ON OrderData (
		BaseCurrency(50),
		QuoteCurrency(50),
		BaseIssuer(50),
		QuoteIssuer(50),
		OrderStatus,
		Network,
		BlockTimeSeconds DESC
	)`)
	a.client.Client.Exec(`CREATE INDEX orderdata_2 ON OrderData (
		Account,
		BaseCurrency(50),
		QuoteCurrency(50),
		BaseIssuer(50),
		QuoteIssuer(50),
		OrderStatus,
		Network,
		BlockTimeSeconds DESC
	)`)
	a.client.Client.Exec(`CREATE INDEX orderdata_3 ON OrderData (
		OrderID
	)`)
	a.client.Client.Exec(`CREATE INDEX orderdata_4 ON OrderData (
		Sequence
	)`)
}
