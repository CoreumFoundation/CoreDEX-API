package http

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/websocket"

	updateproto "github.com/CoreumFoundation/CoreDEX-API/domain/update"
	handler "github.com/CoreumFoundation/CoreDEX-API/utils/httplib/httphandler"
	"github.com/CoreumFoundation/CoreDEX-API/utils/logger"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 4096,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func (s *httpServer) wsEndpoint() handler.Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		// upgrade this connection to a WebSocket connection
		ws, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			logger.Errorf("Error upgrading connection to websocket %+v: %+v", ws, err)
			return err
		}
		s.app.AddSocket(ws)
		err = ws.WriteMessage(1, []byte("Connected"))
		if err != nil {
			logger.Errorf("Error acknowledging websocket connection %+v: %+v", ws, err)
			return err
		}
		// listen indefinitely for new messages coming through on the WebSocket
		s.reader(ws)
		return nil
	}
}

func (s *httpServer) reader(ws *websocket.Conn) {
	for {
		// read a message
		_, p, err := ws.ReadMessage()
		if err != nil {
			if s.app.IsClosed(ws, err) {
				logger.Warnf("reader: Invalid message. Socket is closing? ws: %s, %s: %+v", ws.LocalAddr().String(), ws.RemoteAddr().String(), err)
			}
			return
		}
		m := &updateproto.Subscribe{}
		err = json.Unmarshal(p, m)
		if err != nil {
			logger.Errorf("invalid json %v: %v", p, err)
			continue
		}
		switch m.Action {
		case updateproto.Action_SUBSCRIBE:
			s.app.Subscribe(ws, m)
		case updateproto.Action_UNSUBSCRIBE:
			s.app.Unsubscribe(ws, m)
		case updateproto.Action_CLOSE:
			s.app.Close(ws)
		}
	}
}
