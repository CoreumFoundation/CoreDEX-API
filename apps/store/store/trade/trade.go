package trade

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"

	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/CoreumFoundation/CoreDEX-API/domain/denom"
	"github.com/CoreumFoundation/CoreDEX-API/domain/metadata"
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
Network,
Enriched,
Inverted`

	tradePairTableFields = `Denom1,
Denom2,
MetaData,
PriceTick,
QuantityStep `
)

type Application struct {
	client store.StoreBase
}

// A cache purely a check to see if the value is already in the db (skips a repeated write),
// set is small enough to stay in memory indefinitely
// Reduces trade related writes with 50%
var tradePairCache = make(map[string]bool)

func NewApplication(client *store.StoreBase) *Application {
	app := &Application{
		client: *client,
	}
	app.schema()
	app.index()
	return app
}

// Alphabetical order of the denoms by currency and issuer
func (a *Application) denomInversion(in *tradegrpc.Trade) ([]byte, []byte, string, string, bool, error) {
	den1, den2, inverted := a.denomInverted(in.Denom1, in.Denom2)
	denRet1, err := json.Marshal(den1)
	if err != nil {
		logger.Errorf("Error marshalling denom1 for trade %s-%d-%s: %v", in.TXID, in.Sequence, in.MetaData.Network.String(), err)
		return nil, nil, "", "", false, err
	}
	denRet2, err := json.Marshal(den2)
	if err != nil {
		logger.Errorf("Error marshalling denom2 for trade %s-%d-%s: %v", in.TXID, in.Sequence, in.MetaData.Network.String(), err)
		return nil, nil, "", "", false, err
	}
	return denRet1, denRet2, den1.Denom, den2.Denom, inverted, nil
}

func (*Application) denomInverted(denom1, denom2 *denom.Denom) (*denom.Denom, *denom.Denom, bool) {
	den1 := *denom1
	den2 := *denom2
	inverted := false
	if strings.Compare(den1.Denom, den2.Denom) > 0 {
		den1, den2 = den2, den1
		inverted = true
	}
	return &den1, &den2, inverted
}

func (a *Application) Upsert(in *tradegrpc.Trade) error {
	// Marshal JSON fields
	amount, err := json.Marshal(in.Amount)
	if err != nil {
		logger.Errorf("Error marshalling amount for trade %s-%d-%s: %v", in.TXID, in.Sequence, in.MetaData.Network.String(), err)
		return err
	}
	// Check the symbol order and invert if necessary:
	denom1, denom2, denStr1, denStr2, inverted, err := a.denomInversion(in)
	if err != nil {
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
			    ?, ?, ? ,?, ?,
				? ) 
        ON DUPLICATE KEY UPDATE 
		Amount=?, 
		Price=?, 
		MetaData=?, 
		USD=?,
		Enriched=?`,
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
		in.Enriched,
		inverted,

		amount,
		in.Price,
		metaData,
		in.USD,
		in.Enriched)
	if err != nil {
		logger.Errorf("Error upserting trade %s-%d-%d-%s: %v", in.TXID, in.BlockHeight, in.Sequence, in.MetaData.Network.String(), err)
		return err
	}
	// Reduce the number of writes to the trade pairs table by caching existence of the pairs in memory:
	tradePairKey := a.tradePairKey(denStr1, denStr2, in.MetaData.Network)
	if _, ok := tradePairCache[tradePairKey]; !ok {
		// Keep the trade pairs up to date (ignore the errors: Would only occur on duplicate key or non-recoverable anyway)
		a.client.Client.Exec(`INSERT INTO TradePairs (`+tradePairTableFields+`)
		VALUES (?, ?, ?, NULL, ?)`, denom1, denom2, metaData, 0)
		tradePairCache[tradePairKey] = true
	}
	// And the inverted pair as well
	tradePairKey = a.tradePairKey(denStr2, denStr1, in.MetaData.Network)
	if _, ok := tradePairCache[tradePairKey]; !ok {
		// Keep the trade pairs up to date (ignore the errors: Would only occur on duplicate key or non-recoverable anyway)
		a.client.Client.Exec(`INSERT INTO TradePairs (`+tradePairTableFields+`)
		VALUES (?, ?, ?, NULL, ?)`, denom2, denom1, metaData, 0)
		tradePairCache[tradePairKey] = true
	}
	return nil
}

func (*Application) tradePairKey(denom1, denom2 string, network metadata.Network) string {
	return fmt.Sprintf("%s-%s-%d", denom1, denom2, network)
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
		queryBuilder.WriteString(" AND BlockTimeSeconds >= ?")
		args = append(args, filter.From.AsTime().Unix())
	}
	if filter.To != nil && filter.To.AsTime().Unix() > 0 {
		queryBuilder.WriteString(" AND BlockTimeSeconds < ?")
		args = append(args, filter.To.AsTime().Unix())
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
	// Trades are stored always in the same denom order:
	// Get the denoms in the correct order for the query
	denom1, denom2, _ := a.denomInverted(filter.Denom1, filter.Denom2)
	if denom1 != nil {
		if denom1.Denom != "" {
			queryBuilder.WriteString(" AND Symbol1 = ?")
			args = append(args, denom1.Denom)
		}
	}
	if denom2 != nil {
		if denom2.Denom != "" {
			queryBuilder.WriteString(" AND Symbol2 = ?")
			args = append(args, denom2.Denom)
		}
	}
	if filter.Side != nil {
		queryBuilder.WriteString(" AND Side = ?")
		args = append(args, *filter.Side)
	}
	queryBuilder.WriteString(" ORDER BY BlockTimeSeconds DESC")
	if filter.From == nil || filter.From.AsTime().Unix() == 0 {
		queryBuilder.WriteString(" LIMIT 50")
	}

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
		&trade.Enriched,
		&trade.Inverted,
	)
	if err != nil {
		return nil, err
	}
	json.Unmarshal(amount, &trade.Amount)
	json.Unmarshal(denom1, &trade.Denom1)
	json.Unmarshal(denom2, &trade.Denom2)
	json.Unmarshal(blockTime, &trade.BlockTime)
	json.Unmarshal(metaData, &trade.MetaData)

	// Uninvert the trade if necessary
	if trade.Inverted {
		trade.Denom1, trade.Denom2 = trade.Denom2, trade.Denom1
		trade.Inverted = false
	}
	return trade, nil
}

func (a *Application) GetTradePairs(filter *tradegrpc.TradePairFilter) (*tradegrpc.TradePairs, error) {
	var queryBuilder strings.Builder
	var args []interface{}
	var limit = 1000

	queryBuilder.WriteString(`
			SELECT ` + tradePairTableFields + `
			FROM TradePairs 
			WHERE Network=?
		`)
	args = append(args, filter.Network)

	if filter.Denom1 != nil {
		if filter.Denom1.Currency != "" {
			queryBuilder.WriteString(" AND Currency1 = ?")
			args = append(args, filter.Denom1.Currency)
		}
		if filter.Denom1.Issuer != "" {
			queryBuilder.WriteString(" AND Issuer1 = ?")
			args = append(args, filter.Denom1.Issuer)
		}
	}
	if filter.Denom2 != nil {
		if filter.Denom2.Currency != "" {
			queryBuilder.WriteString(" AND Currency2 = ?")
			args = append(args, filter.Denom2.Currency)
		}
		if filter.Denom2.Issuer != "" {
			queryBuilder.WriteString(" AND Issuer2 = ?")
			args = append(args, filter.Denom2.Issuer)
		}
	}
	queryBuilder.WriteString(" ORDER BY Currency1, Currency2, Issuer1, Issuer2")
	queryBuilder.WriteString(" LIMIT ?")
	args = append(args, limit+1) // +1 to check if there are more results
	var offset int32 = 0
	if filter.Offset != nil {
		queryBuilder.WriteString(" OFFSET ?")
		args = append(args, *filter.Offset)
		offset = *filter.Offset
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
		var denom1, denom2, priceTick []byte
		var quantityStep int64

		if err := rows.Scan(&denom1, &denom2, &metaData, &priceTick, &quantityStep); err != nil {
			return nil, err
		}

		json.Unmarshal(denom1, &tradePair.Denom1)
		json.Unmarshal(denom2, &tradePair.Denom2)
		json.Unmarshal(metaData, &tradePair.MetaData)
		json.Unmarshal(priceTick, &tradePair.PriceTick)
		tradePair.QuantityStep = &quantityStep

		tradePairs = append(tradePairs, &tradePair)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	if len(tradePairs) > limit {
		tradePairs = tradePairs[:limit]
		offset = offset + int32(limit)
	}
	return &tradegrpc.TradePairs{TradePairs: tradePairs, Offset: &offset}, nil
}

func (a *Application) UpsertTradePair(in *tradegrpc.TradePair) error {
	// Marshal JSON fields
	metaData, err := json.Marshal(in.MetaData)
	if err != nil {
		logger.Errorf("Error marshalling metadata for trade pair %s-%s: %v", in.Denom1.Denom, in.Denom2.Denom, err)
		return err
	}
	den1, err := json.Marshal(in.Denom1)
	if err != nil {
		logger.Errorf("Error marshalling denom1 for trade pair %s-%s: %v", in.Denom1.Denom, in.Denom2.Denom, err)
		return err
	}
	den2, err := json.Marshal(in.Denom2)
	if err != nil {
		logger.Errorf("Error marshalling denom2 for trade pair %s-%s: %v", in.Denom1.Denom, in.Denom2.Denom, err)
		return err
	}
	pt, err := json.Marshal(in.PriceTick)
	if err != nil {
		logger.Warnf("Error marshalling priceTick for trade pair %s-%s: %v", in.Denom1.Denom, in.Denom2.Denom, err)
	}
	_, err = a.client.Client.Exec(`INSERT INTO TradePairs (`+tradePairTableFields+`) 
		VALUES (?, ?, ?, ? ,?) 
		ON DUPLICATE KEY UPDATE 
		MetaData=?, PriceTick=?, QuantityStep=?`,
		den1,
		den2,
		metaData,
		pt,
		in.QuantityStep,
		metaData,
		pt,
		in.QuantityStep)
	if err != nil {
		logger.Errorf("Error upserting trade pair %s-%s: %v", in.Denom1.Denom, in.Denom2.Denom, err)
		return err
	}
	return nil
}
