package order

import (
	"database/sql"
	"encoding/json"
	"strings"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"

	ohlcgrpc "github.com/CoreumFoundation/CoreDEX-API/domain/ohlc"
	"github.com/CoreumFoundation/CoreDEX-API/utils/logger"
	store "github.com/CoreumFoundation/CoreDEX-API/utils/mysqlstore"
)

const OHLCDataFields = `Symbol, 
Timestamp, 
Open, 
High, 
Low, 
Close, 
Volume,
QuoteVolume,
NumberOfTrades, 
Period,
PeriodStr,
USDValue, 
MetaData, 
OpenTime, 
CloseTime `

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
// Initialize tables and indexes
func (a *Application) initDB() {
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

func (a *Application) Upsert(in *ohlcgrpc.OHLC) error {
	// Marshal JSON fields
	metaData, err := json.Marshal(in.MetaData)
	if err != nil {
		logger.Errorf("Error marshalling metadata for OHLC %s-%d: %v", in.Symbol, in.Timestamp.AsTime().Unix(), err)
		return err
	}

	period, err := json.Marshal(in.Period)
	if err != nil {
		logger.Errorf("Error marshalling period for OHLC %s-%d: %v", in.Symbol, in.Timestamp.AsTime().Unix(), err)
		return err
	}
	periodStr := in.Period.ToString()
	// Use the mysql client to insert the provided data into the table OHLC
	_, err = a.client.Client.Exec(`INSERT INTO OHLC ( `+OHLCDataFields+` 
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?) 
	ON DUPLICATE KEY UPDATE 
		Open=?, 
		High=?, 
		Low=?, 
		Close=?, 
		Volume=?,
		QuoteVolume=?,
		NumberOfTrades=?, 
		USDValue=?, 
		MetaData=?, 
		OpenTime=?, 
		CloseTime=?`,
		in.Symbol,
		in.Timestamp.AsTime(),
		in.Open,
		in.High,
		in.Low,
		in.Close,
		in.Volume,
		in.QuoteVolume,
		in.NumberOfTrades,
		period,
		periodStr,
		in.USDValue,
		metaData,
		in.OpenTime.AsTime(),
		in.CloseTime.AsTime(),
		in.Open,
		in.High,
		in.Low,
		in.Close,
		in.Volume,
		in.QuoteVolume,
		in.NumberOfTrades,
		in.USDValue,
		metaData,
		in.OpenTime.AsTime(),
		in.CloseTime.AsTime())
	if err != nil {
		logger.Errorf("Error upserting OHLC %s-%d: %v", in.Symbol, in.Timestamp.AsTime().Unix(), err)
		return err
	}
	return nil
}

func (a *Application) Get(filter *ohlcgrpc.OHLCFilter) (*ohlcgrpc.OHLCs, error) {
	ohlcs, err := a.get(filter, false)
	if err != nil {
		return nil, err
	}

	// backfill looks to see if the first datapoint in the request array is allocated,
	// and if not it searches for the previous datapoint
	if filter.Backfill {
		if len(ohlcs) == 0 || (len(ohlcs) > 0 && ohlcs[0].Timestamp.AsTime().Unix() <= filter.From.AsTime().Unix()) {
			backfillOHLCs, err := a.get(filter, true)
			if err != nil {
				return nil, err
			}
			if len(backfillOHLCs) > 0 {
				ohlcs = append(backfillOHLCs, ohlcs...)
			}
		}
	}
	return &ohlcgrpc.OHLCs{OHLCs: ohlcs}, nil
}

func (a *Application) get(filter *ohlcgrpc.OHLCFilter, backFill bool) ([]*ohlcgrpc.OHLC, error) {
	var queryBuilder strings.Builder
	var args []interface{}

	queryBuilder.WriteString(`
		SELECT ` + OHLCDataFields + `
		FROM OHLC 
		WHERE Network=?
	`)
	args = append(args, filter.Network)

	if filter.Symbol != "" {
		queryBuilder.WriteString(" AND Symbol=?")
		args = append(args, filter.Symbol)
	}
	if filter.From != nil && filter.From.AsTime().Unix() > 0 {
		if !backFill && (filter.SingleBucket == nil || !*filter.SingleBucket) {
			queryBuilder.WriteString(" AND Timestamp >= ?")
			args = append(args, filter.From.AsTime().Format(time.DateTime))
		}
	}
	if filter.To != nil && filter.To.AsTime().Unix() > 0 {
		queryBuilder.WriteString(" AND Timestamp < ?")
		args = append(args, filter.To.AsTime().Format(time.DateTime))
	}
	if filter.Period != nil && filter.Period.PeriodType != ohlcgrpc.PeriodType_PERIOD_TYPE_DO_NOT_USE {
		queryBuilder.WriteString(" AND PeriodType=? AND Duration=?")
		args = append(args, filter.Period.PeriodType)
		args = append(args, filter.Period.Duration)
	}
	if !backFill {
		queryBuilder.WriteString(" ORDER BY Timestamp ASC")
	} else {
		// Slightly different query for backfill: There we just want the last datapoint
		queryBuilder.WriteString(" ORDER BY Timestamp DESC LIMIT 1")
	}

	logger.Infof("Query: %s", queryBuilder.String())
	// output the args to the log:
	for i, arg := range args {
		logger.Infof("Arg %d: %v", i, arg)
	}
	rows, err := a.client.Client.Query(queryBuilder.String(), args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ohlcs []*ohlcgrpc.OHLC

	for rows.Next() {
		ohlc, err := mapToOHLC(rows)
		if err != nil {
			return nil, err
		}
		ohlcs = append(ohlcs, ohlc)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return ohlcs, nil
}

func (a *Application) BatchUpsert(in *ohlcgrpc.OHLCs) error {
	for _, ohlc := range in.OHLCs {
		err := a.Upsert(ohlc)
		if err != nil {
			return err
		}
	}
	return nil
}

func (a *Application) GetOHLCsForPeriods(filter *ohlcgrpc.PeriodsFilter) (*ohlcgrpc.OHLCs, error) {
	var queryBuilder strings.Builder
	var args []interface{}

	queryBuilder.WriteString(`
			SELECT ` + OHLCDataFields + `
			FROM OHLC 
			WHERE Symbol=?
		`)
	args = append(args, filter.Symbol)

	if len(filter.Periods) > 0 {
		for i, period := range filter.Periods {
			if i == 0 {
				queryBuilder.WriteString(" AND ( PeriodStr = ? AND Timestamp = ? )")
			}
			if i > 0 {
				queryBuilder.WriteString(" OR ( PeriodStr = ? AND Timestamp = ? )")
			}
			args = append(args, period.Period.ToString())
			args = append(args, period.Timestamp.AsTime())
		}
	}
	rows, err := a.client.Client.Query(queryBuilder.String(), args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	ohlcs := make([]*ohlcgrpc.OHLC, 0)
	for rows.Next() {
		r, err := mapToOHLC(rows)
		if err != nil {
			return nil, err
		}
		ohlcs = append(ohlcs, r)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return &ohlcgrpc.OHLCs{OHLCs: ohlcs}, nil
}

func stringToDate(date string) *time.Time {
	parsed, err := time.Parse("2006-01-02 15:04:05", date)
	if err != nil {
		panic(err)
	}
	return &parsed
}

func mapToOHLC(rows *sql.Rows) (*ohlcgrpc.OHLC, error) {
	var ohlc ohlcgrpc.OHLC
	var timestamp, openTime, closeTime string
	var metaData, period []byte
	var periodStr string // Part of fields for querying, however (by design) not in the OHLC struct
	var quoteVolume sql.NullFloat64

	err := rows.Scan(
		&ohlc.Symbol,
		&timestamp,
		&ohlc.Open,
		&ohlc.High,
		&ohlc.Low,
		&ohlc.Close,
		&ohlc.Volume,
		&quoteVolume,
		&ohlc.NumberOfTrades,
		&period,
		&periodStr,
		&ohlc.USDValue,
		&metaData,
		&openTime,
		&closeTime,
	)
	if err != nil {
		return nil, err
	}

	ohlc.Timestamp = timestamppb.New(*stringToDate(timestamp))
	ohlc.OpenTime = timestamppb.New(*stringToDate(openTime))
	ohlc.CloseTime = timestamppb.New(*stringToDate(closeTime))

	if err := json.Unmarshal(metaData, &ohlc.MetaData); err != nil {
		return nil, err
	}
	if err := json.Unmarshal(period, &ohlc.Period); err != nil {
		return nil, err
	}
	if quoteVolume.Valid {
		ohlc.QuoteVolume = quoteVolume.Float64
	}

	return &ohlc, nil
}
