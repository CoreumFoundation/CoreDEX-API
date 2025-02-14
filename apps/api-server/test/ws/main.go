package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	nethttp "net/http"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/websocket"

	dmn "github.com/CoreumFoundation/CoreDEX-API/apps/api-server/domain"
	"github.com/CoreumFoundation/CoreDEX-API/coreum"
	"github.com/CoreumFoundation/CoreDEX-API/domain/metadata"
	updateproto "github.com/CoreumFoundation/CoreDEX-API/domain/update"
)

func main() {
	c, err := dialSocket()
	if err != nil {
		log.Fatalf("Error dialing socket: %v", err)
	}

	for {
		fmt.Println("Select an option:")
		fmt.Println("0) Exit")
		fmt.Println("1) Test Order book for symbol")
		fmt.Println("2) Test Order book for symbol and account")
		fmt.Println("3) Test Ticker subscription")
		fmt.Println("4) Test OHLC subscription")
		fmt.Println("5) Test Trades for Symbol")
		fmt.Println("6) Test Trades for Account")
		fmt.Println("7) Test Trades for Account and Symbol")
		fmt.Println("8) Events stream (select one or more options first) (Use CTRL+C to exit, or wait 100s)")

		reader := bufio.NewReader(os.Stdin)
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading input:", err)
			continue
		}

		choice, err := strconv.Atoi(input[:len(input)-1])
		if err != nil {
			fmt.Println("Invalid input. Please enter a number between 1 and 7.")
			continue
		}

		switch choice {
		case 0:
			fmt.Println("Exiting...")
			return
		case 1:
			testOrderbookSubscription(c)
		case 2:
			testOrderbookForSymbolAndAccountSubscription(c)
		case 3:
			testTickerSubscription(c)
		case 4:
			testOHLCSubscription(c)
		case 5:
			testTradesForSymbol(c)
		case 6:
			testTradesForAccount(c)
		case 7:
			testTradesForAccountAndSymbol(c)
		case 8:
			testEventsStream(c)
		default:
			fmt.Println("Invalid choice. Please enter a number between 1 and 7.")
		}
	}
}

func dialSocket() (*websocket.Conn, error) {
	d := websocket.Dialer{}
	c, dialResp, err := d.Dial("ws://127.0.0.1:8080/api/ws", nil)

	if got, want := dialResp.StatusCode, nethttp.StatusSwitchingProtocols; got != want {
		log.Printf("dialResp.StatusCode = %q, want %q. Error: %v", got, want, err)
	}
	return c, err
}

func testTickerSubscription(c *websocket.Conn) {
	log.Printf("Testing ticker subscription")
	msg := &updateproto.Subscribe{
		Action: updateproto.Action_SUBSCRIBE,
		Subscription: &updateproto.Subscription{
			Method:  updateproto.Method_TICKER,
			ID:      "dextestdenom0-devcore1p0edzyzpazpt68vdrjy20c42lvwsjpvfzahygs_dextestdenom1-devcore1p0edzyzpazpt68vdrjy20c42lvwsjpvfzahygs",
			Network: metadata.Network_DEVNET,
		},
	}
	m := sendToSocket(c, msg)
	// Decode the content:
	var tickerdata dmn.USDTicker
	err := json.Unmarshal([]byte(m.Subscription.Content), &tickerdata)
	if err != nil {
		log.Fatalf("Error unmarshalling content: %v", err)
	}
	if tickerdata.Tickers == nil || tickerdata.USDTickers == nil {
		log.Printf("Tickers are nil")
		return
	}
	log.Printf("Tickers cur-cur: %+v", *tickerdata.Tickers)
	log.Printf("Tickers USD: %+v", *tickerdata.USDTickers)
}

func sendToSocket(c *websocket.Conn, msg *updateproto.Subscribe) *updateproto.Subscribe {
	c.WriteJSON(msg)
	var respBytes []byte
	// Wait for a response on the message:
	m := &updateproto.Subscribe{}
	for {
		messageType, p, err := c.ReadMessage()
		if err != nil {
			log.Fatalf("Error reading message: %v", err)
		}
		if messageType == websocket.TextMessage {
			if string(p) == "Connected" {
				continue
			}
		}
		// Since we opened the socket just once, we now can have multiple messages incoming based on the users actions.
		// We want to get the message related to the subscription we just made so we need to decode the message and check the ID.
		respBytes = p
		m = &updateproto.Subscribe{}
		err = json.Unmarshal(respBytes, &m)
		if err != nil {
			log.Fatalf("Error unmarshalling message: %v", err)
		}
		if m.Subscription.ID == msg.Subscription.ID && m.Subscription.Method == msg.Subscription.Method {
			break
		}
	}
	// Output the message:
	log.Printf("Received message: %s", string(respBytes))
	return m
}

func testOrderbookSubscription(c *websocket.Conn) {
	log.Printf("Testing orderbook subscription")
	msg := &updateproto.Subscribe{
		Action: updateproto.Action_SUBSCRIBE,
		Subscription: &updateproto.Subscription{
			Method:  updateproto.Method_ORDERBOOK,
			ID:      "dextestdenom0-devcore1p0edzyzpazpt68vdrjy20c42lvwsjpvfzahygs_dextestdenom1-devcore1p0edzyzpazpt68vdrjy20c42lvwsjpvfzahygs",
			Network: metadata.Network_DEVNET,
		},
	}
	m := sendToSocket(c, msg)
	// Decode the content:
	var orders coreum.OrderBookOrders
	err := json.Unmarshal([]byte(m.Subscription.Content), &orders)
	if err != nil {
		log.Fatalf("Error unmarshalling content: %v", err)
	}
	for _, order := range orders.Buy {
		log.Printf("Order: %+v", order)
	}
	for _, order := range orders.Sell {
		log.Printf("Order: %+v", order)
	}
}

func testOrderbookForSymbolAndAccountSubscription(c *websocket.Conn) {
	log.Printf("Testing orderbook subscription for symbol and account")
	msg := &updateproto.Subscribe{
		Action: updateproto.Action_SUBSCRIBE,
		Subscription: &updateproto.Subscription{
			Method:  updateproto.Method_ORDERBOOK_FOR_SYMBOL_AND_ACCOUNT,
			ID:      "devcore1fksu90amj2qgydf43dm2qf6m2dl4szjtx6j5q8_dextestdenom0-devcore1p0edzyzpazpt68vdrjy20c42lvwsjpvfzahygs_dextestdenom1-devcore1p0edzyzpazpt68vdrjy20c42lvwsjpvfzahygs",
			Network: metadata.Network_DEVNET,
		},
	}
	m := sendToSocket(c, msg)
	// Decode the content:
	var orders coreum.OrderBookOrders
	err := json.Unmarshal([]byte(m.Subscription.Content), &orders)
	if err != nil {
		log.Fatalf("Error unmarshalling content: %v", err)
	}
	for _, order := range orders.Buy {
		log.Printf("Order: %+v", order)
	}
	for _, order := range orders.Sell {
		log.Printf("Order: %+v", order)
	}
}

func testOHLCSubscription(c *websocket.Conn) {
	log.Printf("Testing OHLC subscription")
	msg := &updateproto.Subscribe{
		Action: updateproto.Action_SUBSCRIBE,
		Subscription: &updateproto.Subscription{
			Method:  updateproto.Method_OHLC,
			ID:      "dextestdenom0-devcore1p0edzyzpazpt68vdrjy20c42lvwsjpvfzahygs_dextestdenom1-devcore1p0edzyzpazpt68vdrjy20c42lvwsjpvfzahygs_1m",
			Network: metadata.Network_DEVNET,
		},
	}
	sendToSocket(c, msg)
}

func testTradesForSymbol(c *websocket.Conn) {
	log.Printf("Testing trades subscription for symbol")
	msg := &updateproto.Subscribe{
		Action: updateproto.Action_SUBSCRIBE,
		Subscription: &updateproto.Subscription{
			Method:  updateproto.Method_TRADES_FOR_SYMBOL,
			ID:      "dextestdenom0-devcore1p0edzyzpazpt68vdrjy20c42lvwsjpvfzahygs_dextestdenom1-devcore1p0edzyzpazpt68vdrjy20c42lvwsjpvfzahygs",
			Network: metadata.Network_DEVNET,
		},
	}
	sendToSocket(c, msg)
}

func testTradesForAccount(c *websocket.Conn) {
	log.Printf("Testing trades subscription for account")
	msg := &updateproto.Subscribe{
		Action: updateproto.Action_SUBSCRIBE,
		Subscription: &updateproto.Subscription{
			Method:  updateproto.Method_TRADES_FOR_ACCOUNT,
			ID:      "devcore1fksu90amj2qgydf43dm2qf6m2dl4szjtx6j5q8",
			Network: metadata.Network_DEVNET,
		},
	}
	sendToSocket(c, msg)
}

func testTradesForAccountAndSymbol(c *websocket.Conn) {
	log.Printf("Testing trades subscription")
	msg := &updateproto.Subscribe{
		Action: updateproto.Action_SUBSCRIBE,
		Subscription: &updateproto.Subscription{
			Method:  updateproto.Method_TRADES_FOR_ACCOUNT_AND_SYMBOL,
			ID:      "devcore1fksu90amj2qgydf43dm2qf6m2dl4szjtx6j5q8_dextestdenom0-devcore1p0edzyzpazpt68vdrjy20c42lvwsjpvfzahygs_dextestdenom1-devcore1p0edzyzpazpt68vdrjy20c42lvwsjpvfzahygs",
			Network: metadata.Network_DEVNET,
		},
	}
	sendToSocket(c, msg)
}

func testEventsStream(c *websocket.Conn) {
	log.Printf("Testing events stream")
	tStart := time.Now()
	for {
		messageType, p, err := c.ReadMessage()
		if err != nil {
			log.Fatalf("Error reading message: %v", err)
		}
		if messageType == websocket.TextMessage {
			if string(p) == "Connected" {
				continue
			}
		}
		log.Printf("Received message: %s", string(p))
		if time.Since(tStart) > 100*time.Second {
			break
		}
	}
}
