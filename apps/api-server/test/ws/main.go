package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	nethttp "net/http"
	"os"
	"strconv"

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
		default:
			fmt.Println("Invalid choice. Please enter a number between 1 and 7.")
		}
	}
}

func dialSocket() (*websocket.Conn, error) {
	d := websocket.Dialer{}
	c, dialResp, err := d.Dial("ws://127.0.0.1:8080/ws", nil)

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
			ID:      "dextestdenom0-devcore1p0edzyzpazpt68vdrjy20c42lvwsjpvfzahygs_dextestdenom9-devcore1p0edzyzpazpt68vdrjy20c42lvwsjpvfzahygs",
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
	log.Printf("Tickers cur-cur: %+v", *tickerdata.Tickers)
	log.Printf("Tickers USD: %+v", *tickerdata.USDTickers)
}

func sendToSocket(c *websocket.Conn, msg *updateproto.Subscribe) *updateproto.Subscribe {
	c.WriteJSON(msg)
	var respBytes []byte
	// Wait for a response on the message:
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
		respBytes = p
		break // Exit the loop
	}
	// Output the message:
	log.Printf("Received message: %s", string(respBytes))
	// Decode the message into the proto:
	m := &updateproto.Subscribe{}
	err := json.Unmarshal(respBytes, &m)
	if err != nil {
		log.Fatalf("Error unmarshalling message: %v", err)
	}
	return m
}

func testOrderbookSubscription(c *websocket.Conn) {
	log.Printf("Testing orderbook subscription")
	msg := &updateproto.Subscribe{
		Action: updateproto.Action_SUBSCRIBE,
		Subscription: &updateproto.Subscription{
			Method:  updateproto.Method_ORDERBOOK,
			ID:      "dextestdenom0-devcore1p0edzyzpazpt68vdrjy20c42lvwsjpvfzahygs_dextestdenom9-devcore1p0edzyzpazpt68vdrjy20c42lvwsjpvfzahygs",
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
			ID:      "devcore1fpdgztw4aepgy8vezs9hx27yqua4fpewygdspc_dextestdenom0-devcore1p0edzyzpazpt68vdrjy20c42lvwsjpvfzahygs_dextestdenom9-devcore1p0edzyzpazpt68vdrjy20c42lvwsjpvfzahygs",
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
			ID:      "dextestdenom0-devcore1p0edzyzpazpt68vdrjy20c42lvwsjpvfzahygs_dextestdenom9-devcore1p0edzyzpazpt68vdrjy20c42lvwsjpvfzahygs_1m",
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
			ID:      "dextestdenom0-devcore1p0edzyzpazpt68vdrjy20c42lvwsjpvfzahygs_dextestdenom9-devcore1p0edzyzpazpt68vdrjy20c42lvwsjpvfzahygs",
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
			ID:      "devcore1fpdgztw4aepgy8vezs9hx27yqua4fpewygdspc",
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
			ID:      "devcore1fpdgztw4aepgy8vezs9hx27yqua4fpewygdspc_dextestdenom0-devcore1p0edzyzpazpt68vdrjy20c42lvwsjpvfzahygs_dextestdenom9-devcore1p0edzyzpazpt68vdrjy20c42lvwsjpvfzahygs",
			Network: metadata.Network_DEVNET,
		},
	}
	sendToSocket(c, msg)
}
