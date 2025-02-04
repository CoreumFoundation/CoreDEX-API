package main

import (
	"context"
	"encoding/json"
	nethttp "net/http"
	"testing"
	"time"

	"github.com/CoreumFoundation/CoreDEX-API/apps/api-server/app"
	"github.com/CoreumFoundation/CoreDEX-API/apps/api-server/ports/http"
	"github.com/CoreumFoundation/CoreDEX-API/coreum"
	"github.com/CoreumFoundation/CoreDEX-API/domain/metadata"
	updateproto "github.com/CoreumFoundation/CoreDEX-API/domain/update"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/require"
)

func TestWebSocket(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	application := app.NewApplication()
	server := http.NewHttpServer(application)
	go server.Start(ctx)

	<-time.After(1 * time.Second)

	d := websocket.Dialer{}
	c, dialResp, err := d.Dial("ws://127.0.0.1:5354/ws", nil)
	require.NoError(t, err)

	if got, want := dialResp.StatusCode, nethttp.StatusSwitchingProtocols; got != want {
		t.Errorf("dialResp.StatusCode = %q, want %q", got, want)
	}

	msg := &updateproto.Subscribe{
		Action: updateproto.Action_SUBSCRIBE,
		Subscription: &updateproto.Subscription{
			Method:  updateproto.Method_ORDERBOOK,
			ID:      "dextestdenom0-devcore1p0edzyzpazpt68vdrjy20c42lvwsjpvfzahygs_dextestdenom9-devcore1p0edzyzpazpt68vdrjy20c42lvwsjpvfzahygs",
			Network: metadata.Network_DEVNET,
		},
	}
	require.NoError(t, c.WriteJSON(msg))

	messageType, p, err := c.ReadMessage()
	require.NoError(t, err)
	require.Equal(t, websocket.TextMessage, messageType)
	require.Equal(t, []byte("Connected"), p)

	require.NoError(t, c.ReadJSON(msg))
	require.Equal(t, updateproto.Action_RESPONSE, msg.Action)

	var orders coreum.OrderBookOrders
	require.NoError(t, json.Unmarshal([]byte(msg.Subscription.Content), &orders))
	require.Equal(t, uint64(1630), orders.Buy[0].Sequence)
}
