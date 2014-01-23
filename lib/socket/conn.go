package socket

import (
	"github.com/gorilla/websocket"
	"github.com/gorilla/mux"
	"net/http"
)

type connection struct {
	// The websocket connection.
	ws *websocket.Conn
	symbol string
	// Buffered channel of outbound messages.
	send chan []byte
}

func (c *connection) reader() {
	for {
		_, message, err := c.ws.ReadMessage()
		if err != nil {
			break
		}
		h.broadcast <- message
	}
	c.ws.Close()
}

func (c *connection) writer() {
	for message := range c.send {
		err := c.ws.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			break
		}
	}
	c.ws.Close()
}

func WsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ws, err := websocket.Upgrade(w, r, nil, 1024, 1024)
	if _, ok := err.(websocket.HandshakeError); ok {
		http.Error(w, "Not a websocket handshake", 400)
		return
	} else if err != nil {
		return
	}
	c := &connection{send: make(chan []byte, 256), ws: ws, symbol: vars["symbol"]}
	h.register <- c
	defer func() { h.unregister <- c }()
	go c.writer()
	c.reader()
}
