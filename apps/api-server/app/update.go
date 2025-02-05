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
	"google.golang.org/protobuf/types/known/timestamppb"

	dmn "github.com/CoreumFoundation/CoreDEX-API/apps/api-server/domain"
	"github.com/CoreumFoundation/CoreDEX-API/domain/denom"
	"github.com/CoreumFoundation/CoreDEX-API/domain/metadata"
	ohlcgrpc "github.com/CoreumFoundation/CoreDEX-API/domain/ohlc"
	"github.com/CoreumFoundation/CoreDEX-API/domain/symbol"
	tradegrpc "github.com/CoreumFoundation/CoreDEX-API/domain/trade"
	updateproto "github.com/CoreumFoundation/CoreDEX-API/domain/update"
	"github.com/CoreumFoundation/CoreDEX-API/utils/logger"
)

const REFRESH_INTERVAL = 10 * time.Second
const WRITE_CHANNEL_SIZE = 25

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
	for {
		select {
		case <-ctx.Done():
			return
		default:
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
				wg.Add(1)
				// Refresh the subscription:
				switch subscription.Method {
				case updateproto.Method_TRADES_FOR_ACCOUNT:
					go app.updateTradesForAccount(ctx, subscription, startOfInterval, endOfInterval, &wg)
				case updateproto.Method_TRADES_FOR_SYMBOL:
					go app.updateTradesForSymbol(ctx, subscription, startOfInterval, endOfInterval, &wg)
				case updateproto.Method_TRADES_FOR_ACCOUNT_AND_SYMBOL:
					go app.updateTradesForAccountAndSymbol(ctx, subscription, startOfInterval, endOfInterval, &wg)
				case updateproto.Method_TICKER:
					go app.updateTicker(ctx, subscription, &wg)
				case updateproto.Method_OHLC:
					go app.updateOHLC(ctx, subscription, startOfInterval, endOfInterval, &wg)
				case updateproto.Method_ORDERBOOK:
					go app.updateOrderbook(ctx, subscription, &wg)
				case updateproto.Method_ORDERBOOK_FOR_SYMBOL_AND_ACCOUNT:
					go app.updateOrderbookForSymbolAndAccount(ctx, subscription, &wg)
				default:
					wg.Done()
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
	tickers := app.Ticker.GetTickers(ctx, &dmn.TickerReadOptions{
		Symbols: []string{subscription.ID},
		Network: subscription.Network,
	})
	b, err := json.Marshal(tickers)
	if err != nil {
		logger.Errorf("Error marshalling orderbook orders: %v", err)
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
	orders, err := app.Order.OrderBookRelevantOrders(subscription.Network, denoms.Denom1.Denom, denoms.Denom2.Denom, 20)
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
	orders, err := app.Order.OrderBookRelevantOrdersForAccount(subscription.Network, denoms.Denom1.Denom, denoms.Denom2.Denom, account, 20)
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
	logger.Infof("Sending %v to listener: %s %s %s %s,%s", m, subscription.Method, subscription.ID, subscription.Network, ws.LocalAddr().String(), ws.RemoteAddr().String())
}
