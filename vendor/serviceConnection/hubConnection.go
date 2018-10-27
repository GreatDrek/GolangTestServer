package serviceConnection

type Hub struct {
	clients map[*client]bool

	register chan *client

	unregister chan *client
}

func NewHub() *Hub {
	return &Hub{
		register:   make(chan *client),
		unregister: make(chan *client),
		clients:    make(map[*client]bool),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case clientHub := <-h.register:
			h.clients[clientHub] = true
		case clientHub := <-h.unregister:
			if _, ok := h.clients[clientHub]; ok {
				delete(h.clients, clientHub)
				close(clientHub.send)
			}
		}
	}
}
