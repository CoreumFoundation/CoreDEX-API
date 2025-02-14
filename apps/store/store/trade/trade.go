package trade

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"

	"google.golang.org/protobuf/types/known/timestamppb"

	tradegrpc "github.com/CoreumFoundation/CoreDEX-API/domain/trade"
	"github.com/CoreumFoundation/CoreDEX-API/utils/logger"
	store "github.com/CoreumFoundation/CoreDEX-API/utils/mysqlstore"
)

const (
	tradeTableFields = `TXID, 
Account, 
OrderID, 
Sequence, 
Amount,
Price,
Denom1, 
Denom2, 
Side, 
BlockTime, 
BlockHeight,
MetaData, 
USD, 
Network `

	tradePairTableFields = `Denom1,
Denom2,
MetaData `
)

type Application struct {
	client store.StoreBase
}

func NewApplication(client *store.StoreBase) *Application {
	app := &Application{
		client: *client,
	}
	app.initDB()
	return app
}

// Initialize tables and indexes
func (a *Application) initDB() {
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

func (a *Application) Upsert(in *tradegrpc.Trade) error {
	// Marshal JSON fields
	amount, err := json.Marshal(in.Amount)
	if err != nil {
		logger.Errorf("Error marshalling amount for trade %s-%d-%s: %v", in.TXID, in.Sequence, in.MetaData.Network.String(), err)
		return err
	}
	denom1, err := json.Marshal(in.Denom1)
	if err != nil {
		logger.Errorf("Error marshalling denom1 for trade %s-%d-%s: %v", in.TXID, in.Sequence, in.MetaData.Network.String(), err)
		return err
	}
	denom2, err := json.Marshal(in.Denom2)
	if err != nil {
		logger.Errorf("Error marshalling denom2 for trade %s-%d-%s: %v", in.TXID, in.Sequence, in.MetaData.Network.String(), err)
		return err
	}
	blockTime, err := json.Marshal(in.BlockTime)
	if err != nil {
		logger.Errorf("Error marshalling blockTime for trade %s-%d-%s: %v", in.TXID, in.Sequence, in.MetaData.Network.String(), err)
		return err
	}
	if in.MetaData.CreatedAt == nil {
		in.MetaData.CreatedAt = timestamppb.Now()
	}
	in.MetaData.UpdatedAt = timestamppb.Now()
	metaData, err := json.Marshal(in.MetaData)
	if err != nil {
		logger.Errorf("Error marshalling metadata for trade %s-%d-%s: %v", in.TXID, in.Sequence, in.MetaData.Network.String(), err)
		return err
	}

	// Use the mysql client to insert the provided data into the table Trade
	_, err = a.client.Client.Exec(`INSERT INTO Trade (`+tradeTableFields+`) 
        VALUES (?, ?, ?, ?, ?,
			    ?, ?, ?, ?, ?,
			    ?, ?, ? ,? ) 
        ON DUPLICATE KEY UPDATE 
		Amount=?, 
		Price=?, 
		MetaData=?, 
		USD=?`,
		in.TXID,
		in.Account,
		in.OrderID,
		in.Sequence,
		amount,
		in.Price,
		denom1,
		denom2,
		in.Side,
		blockTime,
		in.BlockHeight,
		metaData,
		in.USD,
		in.MetaData.Network,

		amount,
		in.Price,
		metaData,
		in.USD)
	if err != nil {
		logger.Errorf("Error upserting trade %s-%d-%d-%s: %v", in.TXID, in.BlockHeight, in.Sequence, in.MetaData.Network.String(), err)
		return err
	}
	// Keep the trade pairs up to date (ignore the errors: Would only occur on duplicate key or non-recoverable anyway)
	a.client.Client.Exec(`INSERT INTO TradePairs (`+tradePairTableFields+`)
		VALUES (?, ?, ?)`, denom1, denom2, metaData)
	return nil
}

// Get a single trade by ID (Network, TXID, Sequence)
func (a *Application) Get(in *tradegrpc.ID) (*tradegrpc.Trade, error) {
	// Use the mysql client to query for the provided data in the table Trade:
	rows, err := a.client.Client.Query(`SELECT `+tradeTableFields+`
    FROM Trade 
    WHERE 
        TXID=? 
        AND Sequence=? 
        AND Network=?`,
		in.TXID,
		in.Sequence,
		in.Network)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Map the result into the tradegrpc.Trade struct:
	trade := &tradegrpc.Trade{}
	// We are querying by unique key so only get a single result
	if rows.Next() {
		trade, err = mapToTrade(rows)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, fmt.Errorf("no trade found with TXID=%s, Sequence=%d, Network=%d", in.TXID, in.Sequence, in.Network)
	}

	return trade, nil
}

func (a *Application) BatchUpsert(in *tradegrpc.Trades) error {
	for _, trade := range in.Trades {
		err := a.Upsert(trade)
		if err != nil {
			return err
		}
	}
	return nil
}

func (a *Application) GetAll(filter *tradegrpc.Filter) (*tradegrpc.Trades, error) {
	var queryBuilder strings.Builder
	var args []interface{}

	queryBuilder.WriteString(`SELECT ` + tradeTableFields + `
			FROM Trade 
			WHERE Network=?
		`)
	args = append(args, filter.Network)
	if filter.From != nil && filter.From.AsTime().Unix() > 0 {
		queryBuilder.WriteString(" AND JSON_UNQUOTE(JSON_EXTRACT(BlockTime, '$.seconds')) >= ?")
		args = append(args, filter.From.AsTime())
	}
	if filter.To != nil && filter.To.AsTime().Unix() > 0 {
		queryBuilder.WriteString(" AND JSON_UNQUOTE(JSON_EXTRACT(BlockTime, '$.seconds')) < ?")
		args = append(args, filter.To.AsTime())
	}
	if filter.Account != nil && *filter.Account != "" {
		queryBuilder.WriteString(" AND Account=?")
		args = append(args, *filter.Account)
	}
	if filter.Sequence != nil && *filter.Sequence != 0 {
		queryBuilder.WriteString(" AND Sequence=?")
		args = append(args, *filter.Sequence)
	}
	if filter.OrderID != nil && *filter.OrderID != "" {
		queryBuilder.WriteString(" AND OrderID=?")
		args = append(args, *filter.OrderID)
	}
	if filter.TXID != nil && *filter.TXID != "" {
		queryBuilder.WriteString(" AND TXID=?")
		args = append(args, *filter.TXID)
	}
	if filter.Denom1 != nil {
		if filter.Denom1.Denom != "" {
			queryBuilder.WriteString(" AND JSON_UNQUOTE(JSON_EXTRACT(Denom1, '$.Denom')) = ?")
			args = append(args, filter.Denom1.Denom)
		}
	}
	if filter.Denom2 != nil {
		if filter.Denom2.Denom != "" {
			queryBuilder.WriteString(" AND JSON_UNQUOTE(JSON_EXTRACT(Denom2, '$.Denom')) = ?")
			args = append(args, filter.Denom2.Denom)
		}
	}
	queryBuilder.WriteString(" ORDER BY JSON_UNQUOTE(JSON_EXTRACT(BlockTime, '$.seconds')) DESC")

	rows, err := a.client.Client.Query(queryBuilder.String(), args...)
	if err != nil {
		logger.Errorf("Error querying trades: %v", err)
		return nil, err
	}
	defer rows.Close()

	var trades []*tradegrpc.Trade

	for rows.Next() {
		trade, err := mapToTrade(rows)
		if err != nil {
			return nil, err
		}
		trades = append(trades, trade)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return &tradegrpc.Trades{Trades: trades}, nil
}

func mapToTrade(b *sql.Rows) (*tradegrpc.Trade, error) {
	trade := &tradegrpc.Trade{}
	amount := make([]byte, 0)
	denom1 := make([]byte, 0)
	denom2 := make([]byte, 0)
	blockTime := make([]byte, 0)
	metaData := make([]byte, 0)
	var network int // To satisfy the scan
	err := b.Scan(
		&trade.TXID,
		&trade.Account,
		&trade.OrderID,
		&trade.Sequence,
		&amount,
		&trade.Price,
		&denom1,
		&denom2,
		&trade.Side,
		&blockTime,
		&trade.BlockHeight,
		&metaData,
		&trade.USD,
		&network,
	)
	if err != nil {
		return nil, err
	}
	json.Unmarshal(amount, &trade.Amount)
	json.Unmarshal(denom1, &trade.Denom1)
	json.Unmarshal(denom2, &trade.Denom2)
	json.Unmarshal(blockTime, &trade.BlockTime)
	json.Unmarshal(metaData, &trade.MetaData)
	return trade, nil
}

func (a *Application) GetTradePairs(filter *tradegrpc.TradePairFilter) (*tradegrpc.TradePairs, error) {
	var queryBuilder strings.Builder
	var args []interface{}

	queryBuilder.WriteString(`
			SELECT 
				MetaData, 
				Denom1, 
				Denom2
			FROM TradePairs 
			WHERE Network=?
		`)
	args = append(args, filter.Network)

	if filter.Denom1 != nil {
		if filter.Denom1.Currency != "" {
			queryBuilder.WriteString(" AND JSON_UNQUOTE(JSON_EXTRACT(Denom1, '$.Currency')) = ?")
			args = append(args, filter.Denom1.Currency)
		}
		if filter.Denom1.Issuer != "" {
			queryBuilder.WriteString(" AND JSON_UNQUOTE(JSON_EXTRACT(Denom1, '$.Issuer')) = ?")
			args = append(args, filter.Denom1.Issuer)
		}
	}
	if filter.Denom2 != nil {
		if filter.Denom2.Currency != "" {
			queryBuilder.WriteString(" AND JSON_UNQUOTE(JSON_EXTRACT(Denom2, '$.Currency')) = ?")
			args = append(args, filter.Denom2.Currency)
		}
		if filter.Denom2.Issuer != "" {
			queryBuilder.WriteString(" AND JSON_UNQUOTE(JSON_EXTRACT(Denom2, '$.Issuer')) = ?")
			args = append(args, filter.Denom2.Issuer)
		}
	}
	rows, err := a.client.Client.Query(queryBuilder.String(), args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tradePairs []*tradegrpc.TradePair
	for rows.Next() {
		metaData := make([]byte, 0)
		var tradePair tradegrpc.TradePair
		var denom1, denom2 []byte

		if err := rows.Scan(&metaData, &denom1, &denom2); err != nil {
			return nil, err
		}

		json.Unmarshal(denom1, &tradePair.Denom1)
		json.Unmarshal(denom2, &tradePair.Denom2)
		json.Unmarshal(metaData, &tradePair.MetaData)

		tradePairs = append(tradePairs, &tradePair)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return &tradegrpc.TradePairs{TradePairs: tradePairs}, nil
}
