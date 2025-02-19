package ohlc

import "github.com/CoreumFoundation/CoreDEX-API/utils/logger"

// Initialize tables and indexes
func (a *Application) schema() {
	// Create OHLC table
	_, err := a.client.Client.Exec(`CREATE TABLE IF NOT EXISTS OHLC (
        Symbol VARCHAR(255),
        Timestamp TIMESTAMP,
        Open DOUBLE,
        High DOUBLE,
        Low DOUBLE,
        Close DOUBLE,
        Volume DOUBLE,
        NumberOfTrades BIGINT,
        Period JSON,
		PeriodStr VARCHAR(255),
        USDValue DOUBLE,
        MetaData JSON,
        OpenTime TIMESTAMP,
        CloseTime TIMESTAMP,
		Network INT AS (JSON_UNQUOTE(JSON_EXTRACT(MetaData, '$.Network'))),
		PeriodType INT AS (JSON_UNQUOTE(JSON_EXTRACT(Period, '$.PeriodType'))),
		Duration INT AS (JSON_UNQUOTE(JSON_EXTRACT(Period, '$.Duration'))),
        UNIQUE KEY (Symbol, PeriodType, Duration, Timestamp)
    )`)
	if err != nil {
		logger.Fatalf("Error creating OHLC table: %v", err)
	}
	a.client.Client.Exec(`ALTER TABLE OHLC ADD COLUMN QuoteVolume DOUBLE`)
}

func (a *Application) index() {
	a.client.Client.Exec(`CREATE INDEX OHLC_1 ON OHLC(Symbol, PeriodStr, Timestamp)`)
}
