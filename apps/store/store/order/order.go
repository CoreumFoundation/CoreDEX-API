package order

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"

	"google.golang.org/protobuf/types/known/timestamppb"

	ordergrpc "github.com/CoreumFoundation/CoreDEX-API/domain/order"
	"github.com/CoreumFoundation/CoreDEX-API/utils/logger"
	store "github.com/CoreumFoundation/CoreDEX-API/utils/mysqlstore"
)

const OrderDataFields = `Account, 
Type, 
OrderID, 
Sequence, 
BaseDenom, 
QuoteDenom, 
Price, 
Quantity,
RemainingQuantity,
Side, 
GoodTil, 
TimeInForce, 
BlockTime,
OrderFee,
MetaData, 
TXID, 
BlockHeight,
OrderStatus,
Network `

type Application struct {
	client store.StoreBase
}

func NewApplication(client *store.StoreBase) *Application {
	app := &Application{
		client: *client,
	}
	app.schema()
	app.index()
	return app
}

func (a *Application) Get(in *ordergrpc.ID) (*ordergrpc.Order, error) {
	// Use the mysql client to query for the provided data in the table Order:
	rows, err := a.client.Client.Query(`
    SELECT `+OrderDataFields+` 
	FROM OrderData
    WHERE 
        Sequence=? 
        AND Network=?`,
		in.Sequence,
		in.Network)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Map the result into the ordergrpc.Order struct:
	order := &ordergrpc.Order{}
	// We are querying by unique key so only get a single result
	if rows.Next() {
		order, err = mapToOrder(rows)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, fmt.Errorf("no order found with Sequence=%d, Network=%d", in.Sequence, in.Network)
	}

	return order, nil
}
func (a *Application) GetAll(filter *ordergrpc.Filter) (*ordergrpc.Orders, error) {
	var queryBuilder strings.Builder
	var args []interface{}

	queryBuilder.WriteString(`
            SELECT ` + OrderDataFields + `
            FROM OrderData 
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
	if filter.Denom1 != nil {
		if filter.Denom1.Currency != "" {
			queryBuilder.WriteString(" AND BaseCurrency = ?")
			args = append(args, filter.Denom1.Currency)
		}
		if filter.Denom1.Issuer != "" {
			queryBuilder.WriteString(" AND BaseIssuer = ?")
			args = append(args, filter.Denom1.Issuer)
		}
	}
	if filter.Denom2 != nil {
		if filter.Denom2.Currency != "" {
			queryBuilder.WriteString(" AND QuoteCurrency = ?")
			args = append(args, filter.Denom2.Currency)
		}
		if filter.Denom2.Issuer != "" {
			queryBuilder.WriteString(" AND QuoteIssuer = ?")
			args = append(args, filter.Denom2.Issuer)
		}
	}
	if filter.Side != nil && *filter.Side != 0 {
		queryBuilder.WriteString(" AND Side=?")
		args = append(args, *filter.Side)
	}
	if filter.OrderStatus != nil && *filter.OrderStatus != 0 {
		queryBuilder.WriteString(" AND OrderStatus=?")
		args = append(args, *filter.OrderStatus)
	}
	queryBuilder.WriteString(" ORDER BY BlockTimeSeconds DESC")
	rows, err := a.client.Client.Query(queryBuilder.String(), args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []*ordergrpc.Order

	for rows.Next() {
		order, err := mapToOrder(rows)
		if err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return &ordergrpc.Orders{Orders: orders}, nil
}

func (a *Application) Upsert(in *ordergrpc.Order) error {
	// Marshal JSON fields
	baseDenom, err := json.Marshal(in.BaseDenom)
	if err != nil {
		logger.Errorf("Error marshalling baseDenom for order %s-%d-%s: %v", in.OrderID, in.Sequence, in.MetaData.Network.String(), err)
		return err
	}
	quoteDenom, err := json.Marshal(in.QuoteDenom)
	if err != nil {
		logger.Errorf("Error marshalling quoteDenom for order %s-%d-%s: %v", in.OrderID, in.Sequence, in.MetaData.Network.String(), err)
		return err
	}
	blockTime, err := json.Marshal(in.BlockTime)
	if err != nil {
		logger.Errorf("Error marshalling blockTime for order %s-%d-%s: %v", in.OrderID, in.Sequence, in.MetaData.Network.String(), err)
		return err
	}
	if in.MetaData.CreatedAt == nil {
		in.MetaData.CreatedAt = timestamppb.Now()
	}
	in.MetaData.UpdatedAt = timestamppb.Now()
	metaData, err := json.Marshal(in.MetaData)
	if err != nil {
		logger.Errorf("Error marshalling metadata for order %s-%d-%s: %v", in.OrderID, in.Sequence, in.MetaData.Network.String(), err)
		return err
	}
	goodTil, err := json.Marshal(in.GoodTil)
	if err != nil {
		logger.Errorf("Error marshalling goodTil for order %s-%d-%s: %v", in.OrderID, in.Sequence, in.MetaData.Network.String(), err)
		return err
	}
	quantity, err := json.Marshal(in.Quantity)
	if err != nil {
		logger.Errorf("Error marshalling quantity for order %s-%d-%s: %v", in.OrderID, in.Sequence, in.MetaData.Network.String(), err)
		return err
	}
	remainingQuantity, err := json.Marshal(in.RemainingQuantity)
	if err != nil {
		logger.Errorf("Error marshalling remainingQuantity for order %s-%d-%s: %v", in.OrderID, in.Sequence, in.MetaData.Network.String(), err)
		return err
	}
	timeInForce, err := json.Marshal(in.TimeInForce)
	if err != nil {
		logger.Errorf("Error marshalling timeInForce for order %s-%d-%s: %v", in.OrderID, in.Sequence, in.MetaData.Network.String(), err)
		return err
	}
	_, err = a.client.Client.Exec(`INSERT INTO OrderData ( `+OrderDataFields+` ) 
        VALUES (?, ?, ?, ?, ?,
			    ?, ?, ?, ?, ?,
				?, ?, ?, ?, ?,
				?, ?, ?, ?) 
        ON DUPLICATE KEY UPDATE Account=?, 
		Price=?, 
		RemainingQuantity=?,
		BlockTime=?, 
		MetaData=?, 
		TXID=?, 
		BlockHeight=?,
		OrderStatus=?,
		OrderFee=?`,
		in.Account,
		in.Type,
		in.OrderID,
		in.Sequence,
		baseDenom,
		quoteDenom,
		in.Price,
		quantity,
		remainingQuantity,
		in.Side,
		goodTil,
		timeInForce,
		blockTime,
		in.OrderFee,
		metaData,
		*in.TXID,
		in.BlockHeight,
		in.OrderStatus,
		in.MetaData.Network,

		in.Account,
		in.Price,
		remainingQuantity,
		blockTime,
		metaData,
		*in.TXID,
		in.BlockHeight,
		in.OrderStatus,
		in.OrderFee)
	if err != nil {
		logger.Errorf("Error upserting order %s-%d-%s: %v", in.OrderID, in.Sequence, in.MetaData.Network.String(), err)
		return err
	}
	return nil
}

func (a *Application) BatchUpsert(orders *ordergrpc.Orders) error {
	for _, order := range orders.Orders {
		err := a.Upsert(order)
		if err != nil {
			return err
		}
	}
	return nil
}

func mapToOrder(b *sql.Rows) (*ordergrpc.Order, error) {
	order := &ordergrpc.Order{}
	baseDenom := make([]byte, 0)
	quoteDenom := make([]byte, 0)
	goodTil := make([]byte, 0)
	blockTime := make([]byte, 0)
	metaData := make([]byte, 0)
	quantity := make([]byte, 0)
	remainingQuantity := make([]byte, 0)
	var orderStatus sql.NullInt64
	var network int // Dummy variable to scan into

	err := b.Scan(
		&order.Account,
		&order.Type,
		&order.OrderID,
		&order.Sequence,
		&baseDenom,
		&quoteDenom,
		&order.Price,
		&quantity,
		&remainingQuantity,
		&order.Side,
		&goodTil,
		&order.TimeInForce,
		&blockTime,
		&order.OrderFee,
		&metaData,
		&order.TXID,
		&order.BlockHeight,
		&orderStatus,
		&network,
	)
	if err != nil {
		return nil, err
	}
	if orderStatus.Valid {
		order.OrderStatus = ordergrpc.OrderStatus(orderStatus.Int64)
	}

	json.Unmarshal(baseDenom, &order.BaseDenom)
	json.Unmarshal(quoteDenom, &order.QuoteDenom)
	json.Unmarshal(goodTil, &order.GoodTil)
	json.Unmarshal(blockTime, &order.BlockTime)
	json.Unmarshal(metaData, &order.MetaData)
	json.Unmarshal(quantity, &order.Quantity)
	json.Unmarshal(remainingQuantity, &order.RemainingQuantity)
	return order, nil
}
