package serviceConnection

import (
	"encoding/json"
	"log"
	"net/http"
	"serviceAutorization"
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

type dataMesage struct {
	RequestType byte   `json:"requestType"`
	Message     []byte `json:"message"`
}

func typeMessage(data *[]byte) (*dataMesage, error) {
	var inputMessage dataMesage

	err := json.Unmarshal(*data, &inputMessage)
	if err != nil {
		return &inputMessage, err
	}

	return &inputMessage, nil
}

func (c *Сlient) Write(typeM byte, data []byte) {
	var requstMessage dataMesage
	requstMessage.RequestType = typeM
	requstMessage.Message = data

	sendMessage, err := json.Marshal(requstMessage)
	if err != nil {
		c.Disconnect()
		return
	}
	c.send <- sendMessage
}

type iClient interface {
	Read(byte, []byte)
	Inicialization(*Сlient)
	ClientDisconnect()
}

func (c *Сlient) Disconnect() {
	c.disconnect()
}

// Client is a middleman between the websocket connection and the hub.
type Сlient struct {
	hub *Hub

	// The websocket connection.
	conn *websocket.Conn

	// Buffered channel of outbound messages.
	send chan []byte

	registr chan byte

	inClient iClient

	regB bool

	Id int
}

func (c *Сlient) readPump() {
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

		datMessage, err := typeMessage(&message)
		if err != nil {
			break
		}

		if c.regB == false {
			c.registr <- 100
			logginData, err := serviceAutorization.Autorization(datMessage.RequestType, datMessage.Message, c.hub.Db)
			if err != nil {
				log.Println(err)
				break
			} else {
				parseNewClient, err := json.Marshal(logginData)
				if err != nil {
					break
				}
				c.Write(101, parseNewClient)
				c.Id = logginData.Id
			}
			c.regB = true
		}

		c.inClient.Read(datMessage.RequestType, datMessage.Message)
	}
}

func (c *Сlient) writePump() {
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

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	client := &Сlient{hub: hub, conn: conn, send: make(chan []byte, 256), inClient: iC}
	client.hub.register <- client

	iC.Inicialization(client)

	go client.writePump()
	go client.readPump()
	go client.waitAutentification()
}

func (c *Сlient) waitAutentification() {
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

func (c *Сlient) disconnect() {
	c.hub.unregister <- c
	c.conn.Close()
}
