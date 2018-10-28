package serviceConnection

type Hub struct {
	clients map[*Сlient]bool

	register chan *Сlient

	unregister chan *Сlient
}

func NewHub() *Hub {
	return &Hub{
		register:   make(chan *Сlient),
		unregister: make(chan *Сlient),
		clients:    make(map[*Сlient]bool),
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
