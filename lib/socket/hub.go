package socket

import "strings"

type hub struct {
	// Registered connections.
	connections map[*connection]bool

	// Inbound messages from the connections.
	broadcast chan []byte

	// Register requests from the connections.
	register chan *connection

	// Unregister requests from connections.
	unregister chan *connection
}

var h = hub{
	broadcast:   make(chan []byte),
	register:    make(chan *connection),
	unregister:  make(chan *connection),
	connections: make(map[*connection]bool),
}

func GetHub() *hub {
	return &h
}

func SendStock(symbol string, price string) {
	h.broadcast <- []byte(symbol+","+ price)
}

func (h *hub) Run() {
	for {
		select {
		case c := <-h.register:
			h.connections[c] = true
		case c := <-h.unregister:
			delete(h.connections, c)
			close(c.send)
		case m := <-h.broadcast:
			split := strings.Split(string(m),",")
			for c := range h.connections {
				if c.symbol == split[0] {
					select {
					case c.send <- []byte(split[1]):
					default:
						delete(h.connections, c)
						close(c.send)
						go c.ws.Close()
					}
				}
			}
		}
	}
}
