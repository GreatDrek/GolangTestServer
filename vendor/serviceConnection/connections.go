package serviceConnection

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var (
	newline = []byte{'\n'}
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type iClient interface {
	Read([]byte)
	Inicialization(*ConnectionClient)
}

type ConnectionClient struct {
	c *client
}

func (connectionClient *ConnectionClient) Write(data []byte) {
	connectionClient.c.send <- data
}

func (connectionClient *ConnectionClient) Disconnect() {
	connectionClient.c.disconnect()
}

// Client is a middleman between the websocket connection and the hub.
type client struct {
	hub *Hub

	// The websocket connection.
	conn *websocket.Conn

	// Buffered channel of outbound messages.
	send chan []byte

	registr chan byte

	inClient iClient

	regB bool
}

func (c *client) readPump() {
	defer func() {
		log.Println("End readPump")
		c.disconnect()
	}()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

		if c.regB == false {
			c.registr <- 100
			c.regB = true
		}

		c.inClient.Read(message)

	}
}

func (c *client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		log.Println("End writePump")
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}

			w.Write(message)

			// Add queued chat messages to the current websocket message.
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// serveWs handles websocket requests from the peer.
func ServeWs(hub *Hub, w http.ResponseWriter, r *http.Request, iC iClient) {
	log.Println("StartServerWs", time.Now().String())

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	client := &client{hub: hub, conn: conn, send: make(chan []byte, 256), inClient: iC}
	client.hub.register <- client

	iC.Inicialization(&ConnectionClient{client})

	go client.writePump()
	go client.readPump()
	go client.waitAutentification()
	log.Println("EndServerWs", time.Now().String())
}

func (c *client) waitAutentification() {
	c.registr = make(chan byte)
	defer func() {
		close(c.registr)
		log.Println("End waitAutentification")
	}()

	timer := time.NewTimer(time.Second * 15)

	select {
	case <-timer.C:
		c.disconnect()
		break
	case b := <-c.registr:
		switch b {
		case 100:
			//log.Println("stop")
			break
		default:
			log.Println("error number")
			c.disconnect()
			break
		}
		break
	}
}

func (c *client) disconnect() {
	c.hub.unregister <- c
	c.conn.Close()
}
