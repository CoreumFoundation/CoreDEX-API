package app

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/samber/lo"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/CoreumFoundation/CoreDEX-API/apps/api-server/app/ticker"
	dmn "github.com/CoreumFoundation/CoreDEX-API/apps/api-server/domain"
	"github.com/CoreumFoundation/CoreDEX-API/domain/denom"
	"github.com/CoreumFoundation/CoreDEX-API/domain/metadata"
	ohlcgrpc "github.com/CoreumFoundation/CoreDEX-API/domain/ohlc"
	orderproperties "github.com/CoreumFoundation/CoreDEX-API/domain/order-properties"
	"github.com/CoreumFoundation/CoreDEX-API/domain/symbol"
	tradegrpc "github.com/CoreumFoundation/CoreDEX-API/domain/trade"
	updateproto "github.com/CoreumFoundation/CoreDEX-API/domain/update"
	"github.com/CoreumFoundation/CoreDEX-API/utils/logger"
)

const (
	REFRESH_INTERVAL   = 1 * time.Second
	WRITE_CHANNEL_SIZE = 25
	OHLC_REFRESH       = 10 // Counter based on the REFRESH_INTERVAL in seconds (so 10 is 10 seconds)
	WALLET_REFRESH     = 10 // Counter based on the REFRESH_INTERVAL in seconds (so 10 is 10 seconds)
)

type Message struct {
	Network metadata.Network
	Method  updateproto.Method
	ID      string
	Content interface{}
}

type queueItem struct {
	ws  *websocket.Conn
	sub *updateproto.Subscription
}

type subscriptionManager struct {
	*updateproto.Subscription
	bufferChan chan queueItem
}

var (
	// Allow subscribing to the different types and methods for a certain websocket with potentially more than 1 key to listen for
	// first map[string] contains the network identifier
	// Map of network, method, websocket, id (id can be empty if it is a broadcast style message)
	listeners         = make(map[*websocket.Conn][]*updateproto.Subscription)
	senders           = make(map[*websocket.Conn]*subscriptionManager)
	senderMutex       = sync.RWMutex{}
	subscriptionMutex = sync.RWMutex{}
)

// StartUpdater Iterates over the listeners and sends the message to the listeners
func (app *Application) StartUpdater(ctx context.Context) {
	currentTime := time.Now()
	refreshCounter := 0
	for {
		select {
		case <-ctx.Done():
			return
		default:
			refreshCounter++
			uniqueSubscriptions := make(map[string]*updateproto.Subscription)
			// Identify unique subscriptions:
			subscriptionMutex.RLock()
			for _, subscriptions := range listeners {
				for _, subscription := range subscriptions {
					uniqueSubscriptions[fmt.Sprintf("%s-%s-%s", subscription.Network, subscription.Method, subscription.ID)] = subscription
				}
			}
			subscriptionMutex.RUnlock()
			// Refresh the unique subscriptions:
			// Synchronize the refresh interval at the modulus of the REFRESH_INTERVAL
			// Calculate the start of the current interval
			endOfInterval := currentTime.Truncate(REFRESH_INTERVAL)
			startOfInterval := endOfInterval.Add(-REFRESH_INTERVAL)
			// note for later:
			// It is unknown how long all the refreshes will take. There can be 100s of subscriptions, leading to congestion.
			// Using go routines and wait groups to refresh the subscriptions concurrently
			// If this does not work out, we can delegate the unique subscriptions (for accounts) to the update process itself
			wg := sync.WaitGroup{}
			for _, subscription := range uniqueSubscriptions {
				// Refresh the subscription:
				switch subscription.Method {
				case updateproto.Method_TRADES_FOR_ACCOUNT:
					wg.Add(1)
					go app.updateTradesForAccount(ctx, subscription, startOfInterval, endOfInterval, &wg)
				case updateproto.Method_TRADES_FOR_SYMBOL:
					wg.Add(1)
					go app.updateTradesForSymbol(ctx, subscription, startOfInterval, endOfInterval, &wg)
				case updateproto.Method_TRADES_FOR_ACCOUNT_AND_SYMBOL:
					wg.Add(1)
					go app.updateTradesForAccountAndSymbol(ctx, subscription, startOfInterval, endOfInterval, &wg)
				case updateproto.Method_TICKER:
					// dynamic refresh interval based on the cache duration TICKER_CACHE
					if refreshCounter%(int(ticker.TICKER_CACHE.Seconds())) == 0 {
						wg.Add(1)
						go app.updateTicker(ctx, subscription, &wg)
					}
				case updateproto.Method_OHLC:
					// Reduced number of updates: Updating the OHLC to aggressive can lead to overload of FE, and to overload of the DB:
					if refreshCounter%OHLC_REFRESH == 0 {
						wg.Add(1)
						soi := currentTime.Add(-OHLC_REFRESH * time.Second)
						/* From (startOfInterval, soi) needs to be corrected for the fact that the OHLC data (trade data) is delayed in the processing compared to the clock:
						A block gets processed every 1.1 seconds, however can be fully for the previous bucket involved.
						The from thus needs to include the previous bucket as well to be certain we have all the data.
						Or more precise: The from needs to include the previous bucket if the current time is within the first 1.1 seconds of the current bucket,
						add some margin for some delays: Using the refresh interval as the margin for the delay should be sufficient.
						This can be done by setting the start of interval before the minute when we are in the first 5 seconds of the minute as indicated by soi
						The underlying interval calculations will truncate the timestamp to the start of the actual interval involved (based on the period)
						*/
						if soi.Second() <= OHLC_REFRESH {
							soi = soi.Add(-time.Minute)
						}
						go app.updateOHLC(ctx, subscription, soi, endOfInterval, &wg)
					}
				case updateproto.Method_ORDERBOOK:
					wg.Add(1)
					go app.updateOrderbook(ctx, subscription, &wg)
				case updateproto.Method_ORDERBOOK_FOR_SYMBOL_AND_ACCOUNT:
					wg.Add(1)
					go app.updateOrderbookForSymbolAndAccount(ctx, subscription, &wg)
				case updateproto.Method_WALLET:
					if refreshCounter%WALLET_REFRESH == 0 {
						wg.Add(1)
						go app.updateWallet(subscription, &wg)
					}
				}
			}
			wg.Wait()
			// Unique subscriptions are refreshed, and the content has been added/updated

			// Send out the data to the listeners:
			for ws, subscriptions := range listeners {
				for _, subscription := range subscriptions {
					if sub, ok := uniqueSubscriptions[fmt.Sprintf("%s-%s-%s", subscription.Network, subscription.Method, subscription.ID)]; ok {
						app.addToQueue(ctx, ws, sub)
					}
				}
			}
			// Set start of next interval avoiding drift:
			currentTime = endOfInterval.Add(REFRESH_INTERVAL)
			// Sleep until the next interval:
			if time.Until(currentTime) > 0 {
				time.Sleep(time.Until(currentTime))
			}
		}
	}
}

/* Self terminating pool of go routines to write to the websockets
The self termination design here has been choosen to avoid having to cleanup go routines when the client disconnects
The process is as follows:
1. Check if there is a go routine for the socket in the senders map, if not:
2. Create a channel for the socket in the senders map
3. Start a go routine for the socket
About the go routine:
4. Go routine has a time out on which it terminates itself
5. go routine removes itself from the list of senders

Since we use the 3*REFRESH_INTERVAL as the time out for termination, the number of initializations and terminations is limited
*/

func (app *Application) addToQueue(ctx context.Context, ws *websocket.Conn, sub *updateproto.Subscription) {
	senderMutex.Lock()
	if _, ok := senders[ws]; !ok {
		senders[ws] = &subscriptionManager{
			Subscription: sub,
			bufferChan:   make(chan queueItem, WRITE_CHANNEL_SIZE),
		}
		go app.startSender(ctx, ws)
	}
	if len(senders[ws].bufferChan) < WRITE_CHANNEL_SIZE {
		senders[ws].bufferChan <- queueItem{ws, sub}
	}
	senderMutex.Unlock()
}

func (app *Application) startSender(ctx context.Context, ws *websocket.Conn) {
	// Detach the channel from the sender map so that other manipulations can be done on the map (error handling can remove the socket from the map)
	senderChan := senders[ws].bufferChan
	for {
		select {
		case <-ctx.Done():
			return
		case writeRequest := <-senderChan:
			app.writeMessage(writeRequest.ws, writeRequest.sub)
		case <-time.After(3 * REFRESH_INTERVAL):
			senderMutex.Lock()
			close(senderChan)
			delete(senders, ws)
			senderMutex.Unlock()
			return
		}
	}
}

func (app *Application) updateTradesForAccount(ctx context.Context, subscription *updateproto.Subscription, from, to time.Time, wg *sync.WaitGroup) {
	trades, err := app.Trade.GetTrades(ctx, &tradegrpc.Filter{
		Network: subscription.Network,
		From:    timestamppb.New(from),
		To:      timestamppb.New(to),
		Account: &subscription.ID,
	})
	if err != nil {
		logger.Errorf("Error getting trades: %v", err)
		wg.Done()
		return
	}
	b, err := json.Marshal(trades)
	if err != nil {
		logger.Errorf("Error marshalling trades: %v", err)
		wg.Done()
		return
	}
	subscription.Content = string(b)
	wg.Done()
}

func (app *Application) updateTradesForSymbol(ctx context.Context, subscription *updateproto.Subscription, from, to time.Time, wg *sync.WaitGroup) {
	// The denoms are concatenated with a separator _
	denoms, err := symbol.NewSymbol(subscription.ID)
	if err != nil {
		logger.Errorf("Error parsing denoms: %v", err)
		wg.Done()
		return
	}
	trades, err := app.Trade.GetTrades(ctx, &tradegrpc.Filter{
		Network: subscription.Network,
		From:    timestamppb.New(from),
		To:      timestamppb.New(to),
		Denom1:  denoms.Denom1,
		Denom2:  denoms.Denom2,
		Side:    lo.ToPtr(orderproperties.Side_SIDE_BUY),
	})
	if err != nil {
		logger.Errorf("Error getting trades: %v", err)
		wg.Done()
		return
	}
	b, err := json.Marshal(trades)
	if err != nil {
		logger.Errorf("Error marshalling trades: %v", err)
		wg.Done()
		return
	}
	subscription.Content = string(b)
	wg.Done()
}

func (app *Application) updateTradesForAccountAndSymbol(ctx context.Context, subscription *updateproto.Subscription, from, to time.Time, wg *sync.WaitGroup) {
	// The account and denoms are concatenated with a separator _
	parts := strings.SplitN(subscription.ID, "_", 2)
	if len(parts) != 2 {
		logger.Errorf("Error parsing ID")
		wg.Done()
		return
	}
	account := parts[0]
	denoms, err := symbol.NewSymbol(parts[1])
	if err != nil {
		logger.Errorf("Error parsing denoms: %v", err)
		wg.Done()
		return
	}
	trades, err := app.Trade.GetTrades(ctx, &tradegrpc.Filter{
		Network: subscription.Network,
		From:    timestamppb.New(from),
		To:      timestamppb.New(to),
		Account: &account,
		Denom1:  denoms.Denom1,
		Denom2:  denoms.Denom2,
	})
	if err != nil {
		logger.Errorf("Error getting trades: %v", err)
		wg.Done()
		return
	}
	b, err := json.Marshal(trades)
	if err != nil {
		logger.Errorf("Error marshalling trades: %v", err)
		wg.Done()
		return
	}
	subscription.Content = string(b)
	wg.Done()
}

func (app *Application) updateTicker(ctx context.Context, subscription *updateproto.Subscription, wg *sync.WaitGroup) {
	opt := dmn.NewTickerReadOptions([]string{subscription.ID}, time.Now().Truncate(time.Second), 24*time.Hour)
	opt.Network = subscription.Network
	tickers := app.Ticker.GetTickers(ctx, opt)
	b, err := json.Marshal(tickers)
	if err != nil {
		logger.Errorf("Error marshalling tickers: %v", err)
		wg.Done()
		return
	}
	subscription.Content = string(b)
	wg.Done()
}

func (app *Application) updateOHLC(ctx context.Context, subscription *updateproto.Subscription, startOfInterval, endOfInterval time.Time, wg *sync.WaitGroup) {
	// The denoms and the period (interval/bucket) are concatenated with a separator _ in the requesting ID (denom-issuer_denom2-issuer2_interval)
	// where interval is the same as in the restful call to the ohlc endpoint
	denomsPeriod := strings.Split(subscription.ID, "_")
	if len(denomsPeriod) != 3 {
		logger.Infof("Error parsing denoms and period (incorrect format): %v", denomsPeriod)
		wg.Done()
		return
	}
	// Parse into the correct format:
	denom1, err := denom.NewDenom(denomsPeriod[0])
	if err != nil {
		logger.Errorf("Error parsing denom1: %v", err)
		wg.Done()
		return
	}
	denom2, err := denom.NewDenom(denomsPeriod[1])
	if err != nil {
		logger.Errorf("Error parsing denom2: %v", err)
		wg.Done()
		return
	}
	period, err := dmn.HttpPeriodToPeriod(denomsPeriod[2])
	if err != nil {
		logger.Errorf("Error parsing interval: %v", err)
		wg.Done()
		return
	}
	from := period.ToOHLCKeyTimestampFrom(startOfInterval.UnixNano())
	to := period.ToOHLCKeyTimestampTo(endOfInterval.UnixNano())
	ohlcs, err := app.OHLC.Get(ctx, &ohlcgrpc.OHLCFilter{
		Symbol:  denom1.Denom + "_" + denom2.Denom,
		Period:  period,
		From:    timestamppb.New(time.Unix(0, from)),
		To:      timestamppb.New(time.Unix(0, to)),
		Network: subscription.Network,
	})
	if err != nil {
		logger.Errorf("Error getting OHLCs: %v", err)
		wg.Done()
		return
	}
	b, err := json.Marshal(ohlcs)
	if err != nil {
		logger.Errorf("Error marshalling OHLCs: %v", err)
		wg.Done()
		return
	}
	subscription.Content = string(b)
	logger.Infof("Filter: %v, %v, %v, %v, %v, %v", subscription.Network, denom1.Denom, denom2.Denom, period, startOfInterval, endOfInterval)
	for _, ohlc := range ohlcs {
		logger.Infof("OHLC: %v, %v, %v, %v, %v, %v", ohlc[0], ohlc[1], ohlc[2], ohlc[3], ohlc[4], ohlc[5])
	}
	wg.Done()
}

func (app *Application) updateOrderbook(_ context.Context, subscription *updateproto.Subscription, wg *sync.WaitGroup) {
	// The account and denoms are concatenated with a separator _
	denoms, err := symbol.NewSymbol(subscription.ID)
	if err != nil {
		logger.Errorf("Error parsing denoms: %v", err)
		wg.Done()
		return
	}
	orders, err := app.Order.OrderBookRelevantOrders(subscription.Network, denoms.Denom1.Denom, denoms.Denom2.Denom, 50, true)
	if err != nil {
		logger.Errorf("Error getting orderbook orders: %v", err)
		wg.Done()
		return
	}
	b, err := json.Marshal(orders)
	if err != nil {
		logger.Errorf("Error marshalling orderbook orders: %v", err)
		wg.Done()
		return
	}
	subscription.Content = string(b)
	wg.Done()
}

func (app *Application) updateOrderbookForSymbolAndAccount(_ context.Context, subscription *updateproto.Subscription, wg *sync.WaitGroup) {
	// The account and denoms are concatenated with a separator _
	parts := strings.SplitN(subscription.ID, "_", 2)
	if len(parts) != 2 {
		logger.Errorf("Error parsing ID")
		wg.Done()
		return
	}
	account := parts[0]
	denoms, err := symbol.NewSymbol(parts[1])
	if err != nil {
		logger.Errorf("Error parsing denoms: %v", err)
		wg.Done()
		return
	}
	orders, err := app.Order.OrderBookRelevantOrdersForAccount(subscription.Network, denoms.Denom1.Denom, denoms.Denom2.Denom, account)
	if err != nil {
		logger.Errorf("Error getting orderbook orders: %v", err)
		wg.Done()
		return
	}
	b, err := json.Marshal(orders)
	if err != nil {
		logger.Errorf("Error marshalling orderbook orders: %v", err)
		wg.Done()
		return
	}
	subscription.Content = string(b)
	wg.Done()
}

func (app *Application) updateWallet(subscription *updateproto.Subscription, wg *sync.WaitGroup) {
	wallet, err := app.Order.WalletAssets(subscription.Network, subscription.ID)
	if err != nil {
		logger.Errorf("Error getting wallet: %v", err)
		wg.Done()
		return
	}
	b, err := json.Marshal(wallet)
	if err != nil {
		logger.Errorf("Error marshalling wallet: %v", err)
		wg.Done()
		return
	}
	subscription.Content = string(b)
	wg.Done()
}

func (*Application) AddSocket(ws *websocket.Conn) {
	subscriptionMutex.RLock()
	listeners[ws] = make([]*updateproto.Subscription, 0)
	subscriptionMutex.RUnlock()
}

// Subscribe the websocket connection for the given method and type.
// Note: The mutex.Unlock is specifically not using defer for performance reasons
// (defer can take 10s to 100s of milliseconds in the go scheduler before it is executed)
func (*Application) Subscribe(ws *websocket.Conn, update *updateproto.Subscribe) {
	subscriptionMutex.Lock()
	// Check if ws is present in the map:
	if listeners[ws] == nil {
		listeners[ws] = make([]*updateproto.Subscription, 0)
	}
	// Add the subscription to the map (Check if the subscription is already present):
	for _, subscription := range listeners[ws] {
		if subscription.Network == update.Subscription.Network &&
			subscription.Method == update.Subscription.Method &&
			subscription.ID == update.Subscription.ID {
			subscriptionMutex.Unlock()
			return
		}
	}
	listeners[ws] = append(listeners[ws], update.Subscription)
	subscriptionMutex.Unlock()
}

// Unsubscribe the websocket connection for the given method and type.
// Note: The mutex.Unlock is specifically not using defer for performance reasons
// (defer can take 10s to 100s of milliseconds in the go scheduler before it is executed)
func (*Application) Unsubscribe(ws *websocket.Conn, update *updateproto.Subscribe) {
	subscriptionMutex.Lock()
	// Check if ws is present in the map:
	if listeners[ws] == nil {
		subscriptionMutex.Unlock()
		return
	}
	deleteSub := make([]*updateproto.Subscription, 0)
	// Delete the subscription from the map:
	for _, subscription := range listeners[ws] {
		if subscription.Network == update.Subscription.Network &&
			subscription.Method == update.Subscription.Method &&
			subscription.ID == update.Subscription.ID {
			deleteSub = append(deleteSub, subscription)
		}
	}
	// Use the deleteSub array to remove the subscriptions from the listeners map
	for _, subscription := range deleteSub {
		for i, sub := range listeners[ws] {
			if sub == subscription {
				listeners[ws] = append(listeners[ws][:i], listeners[ws][i+1:]...)
				break
			}
		}
	}
	subscriptionMutex.Unlock()
}

// Close is called by both the client and by the sub processes trying to communicate with the client:
// If the connection is gone, the sub process will call the close so that other sub processes don't have to try
// and clean up the connection individually.
func (*Application) Close(ws *websocket.Conn) {
	subscriptionMutex.Lock()
	// Check if ws is present in the map:
	if listeners[ws] == nil {
		subscriptionMutex.Unlock()
		return
	}
	// Delete the subscription from the map:
	delete(listeners, ws)
	subscriptionMutex.Unlock()
	ws.Close()
}

func (*Application) IsClosed(ws *websocket.Conn, wsErr error) bool {
	switch {
	case wsErr == nil:
		return false
	case strings.Contains(wsErr.Error(), "1005"):
		logger.Infof("Close 1005 for websocket %s,%s: %s", ws.LocalAddr().String(), ws.RemoteAddr().String(), wsErr.Error())
		return true
	case errors.Is(wsErr, websocket.ErrCloseSent):
		logger.Infof("Close send by websocket %s,%s: %s", ws.LocalAddr().String(), ws.RemoteAddr().String(), wsErr.Error())
		return true
	default:
		logger.Warnf("Unknown error writing to websocket %s,%s: %s", ws.LocalAddr().String(), ws.RemoteAddr().String(), wsErr.Error())
		return true
	}

}

func (app *Application) writeMessage(ws *websocket.Conn, subscription *updateproto.Subscription) {
	m := &updateproto.Subscribe{
		Action: updateproto.Action_RESPONSE,
		Subscription: &updateproto.Subscription{
			Method:  subscription.Method,
			ID:      subscription.ID,
			Network: subscription.Network,
			Content: subscription.Content,
		},
	}
	if app.IsClosed(ws, ws.WriteJSON(m)) {
		// The connection is dead, remove from the map
		app.Close(ws)
		return
	}
}
