package currency

import (
	"database/sql"
	"encoding/json"

	currencygrpc "github.com/CoreumFoundation/CoreDEX-API/domain/currency"
	"github.com/CoreumFoundation/CoreDEX-API/domain/decimal"
	"github.com/CoreumFoundation/CoreDEX-API/domain/denom"
	"github.com/CoreumFoundation/CoreDEX-API/domain/metadata"
	"github.com/CoreumFoundation/CoreDEX-API/utils/logger"
	store "github.com/CoreumFoundation/CoreDEX-API/utils/mysqlstore"
)

const currencyTableFields = `Denom, 
SendCommission, 
BurnRate, 
InitialAmount, 
Chain,
OriginChain, 
ChainSupply, 
Description, 
SkipDisplay, 
MetaData, 
Network,
DenomString `

type Application struct {
	client store.StoreBase
}

func NewApplication(client *store.StoreBase) *Application {
	app := &Application{
		client: *client,
	}
	app.schema()
	return app
}

func (a *Application) Upsert(in *currencygrpc.Currency) error {
	// Marshal JSON fields
	denom, err := json.Marshal(in.Denom)
	if err != nil {
		logger.Errorf("Error marshalling denom for currency on chain %s: %v", in.Chain, err)
		return err
	}
	metaData, err := json.Marshal(in.MetaData)
	if err != nil {
		logger.Errorf("Error marshalling metadata for currency on chain %s: %v", in.Chain, err)
		return err
	}
	sendCommission, err := json.Marshal(in.SendCommission)
	if err != nil {
		logger.Errorf("Error marshalling sendCommission for currency on chain %s: %v", in.Chain, err)
		return err
	}
	burnRate, err := json.Marshal(in.BurnRate)
	if err != nil {
		logger.Errorf("Error marshalling burnRate for currency on chain %s: %v", in.Chain, err)
		return err
	}
	initialAmount, err := json.Marshal(in.InitialAmount)
	if err != nil {
		logger.Errorf("Error marshalling initialAmount for currency on chain %s: %v", in.Chain, err)
		return err
	}

	// Use the mysql client to insert the provided data into the table Currency
	_, err = a.client.Client.Exec(`INSERT IGNORE INTO Currency ( `+currencyTableFields+`
    ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		denom,
		sendCommission,
		burnRate,
		initialAmount,
		in.Chain,
		in.OriginChain,
		in.ChainSupply,
		in.Description,
		in.SkipDisplay,
		metaData,
		in.MetaData.Network,
		in.Denom.Denom)
	if err != nil {
		logger.Errorf("Error upserting currency on chain %s: %v", in.Chain, err)
		return err
	}
	return nil
}

func (a *Application) Get(id *currencygrpc.ID) (*currencygrpc.Currency, error) {
	// Query the database for the currency record
	rows, err := a.client.Client.Query(`SELECT `+currencyTableFields+`
        FROM Currency 
        WHERE Network = ? AND DenomString = ?`, id.Network, id.Denom)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // No record found
		}
		logger.Errorf("Error retrieving currency for network %s and denom %s: %v", id.Network.String(), id.Denom, err)
		return nil, err
	}
	t := rows.Next()
	defer rows.Close()
	if !t {
		return nil, nil
	}
	return mapToCurrency(rows)
}

func (a *Application) BatchUpsert(in *currencygrpc.Currencies) error {
	for _, cur := range in.Currencies {
		err := a.Upsert(cur)
		if err != nil {
			return err
		}
	}
	return nil
}

func mapToCurrency(b *sql.Rows) (*currencygrpc.Currency, error) {
	var (
		denomb         []byte
		sendCommission []byte
		burnRate       []byte
		initialAmount  []byte
		chain          string
		chainSupply    string
		originChain    string
		description    string
		skipDisplay    bool
		metaData       []byte
		network        int
		denomString    string
	)

	err := b.Scan(
		&denomb,
		&sendCommission,
		&burnRate,
		&initialAmount,
		&chain,
		&originChain,
		&chainSupply,
		&description,
		&skipDisplay,
		&metaData,
		&network,
		&denomString,
	)
	if err != nil {
		logger.Errorf("Error scanning currency: %v", err)
		return nil, err
	}
	// Unmarshal JSON fields
	var denomStruct denom.Denom
	err = json.Unmarshal(denomb, &denomStruct)
	if err != nil {
		logger.Errorf("Error unmarshalling denom for currency: %v", err)
		return nil, err
	}

	var metaDataStruct metadata.MetaData
	err = json.Unmarshal(metaData, &metaDataStruct)
	if err != nil {
		logger.Errorf("Error unmarshalling metadata for currency %s: %v", err)
		return nil, err
	}

	var sendCommissionStruct decimal.Decimal
	err = json.Unmarshal(sendCommission, &sendCommissionStruct)
	if err != nil {
		logger.Errorf("Error unmarshalling sendCommission for currency: %v", err)
		return nil, err
	}

	var burnRateStruct decimal.Decimal
	err = json.Unmarshal(burnRate, &burnRateStruct)
	if err != nil {
		logger.Errorf("Error unmarshalling burnRate for currency: %v", err)
		return nil, err
	}

	var initialAmountStruct decimal.Decimal
	err = json.Unmarshal(initialAmount, &initialAmountStruct)
	if err != nil {
		logger.Errorf("Error unmarshalling initialAmount for currency: %v", err)
		return nil, err
	}

	// Create and return the Currency struct
	currency := &currencygrpc.Currency{
		Denom:          &denomStruct,
		SendCommission: &sendCommissionStruct,
		BurnRate:       &burnRateStruct,
		InitialAmount:  &initialAmountStruct,
		Chain:          chain,
		OriginChain:    originChain,
		ChainSupply:    chainSupply,
		Description:    description,
		SkipDisplay:    skipDisplay,
		MetaData:       &metaDataStruct,
	}

	return currency, nil
}

func (a *Application) GetAll(filter *currencygrpc.Filter) (*currencygrpc.Currencies, error) {
	var (
		query = `SELECT ` + currencyTableFields + ` FROM Currency WHERE Network = ?`
		args  = []interface{}{filter.Network}
	)

	// Add optional Denom filter
	if filter.Denom != nil {
		query += ` AND JSON_UNQUOTE(JSON_EXTRACT(Denom, '$.Denom')) = ?`
		args = append(args, filter.Denom.Denom)
	}

	// Execute the query
	rows, err := a.client.Client.Query(query, args...)
	if err != nil {
		logger.Errorf("Error querying currencies: %v", err)
		return nil, err
	}
	defer rows.Close()

	var currencies []*currencygrpc.Currency

	for rows.Next() {
		currency, err := mapToCurrency(rows)
		if err != nil {
			logger.Errorf("Error mapping currency: %v", err)
			return nil, err
		}
		currencies = append(currencies, currency)
	}

	if err = rows.Err(); err != nil {
		logger.Errorf("Error iterating currency rows: %v", err)
		return nil, err
	}

	return &currencygrpc.Currencies{Currencies: currencies}, nil
}
