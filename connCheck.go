// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/json"
	"io/ioutil"

	"bytes"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
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
	newline           = []byte{'\n'}
	googleOauthConfig *oauth2.Config
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	hub *Hub

	// The websocket connection.
	conn *websocket.Conn

	// Buffered channel of outbound messages.
	send chan []byte

	registr chan byte

	logginData *logginDataClient
}

type logginDataClient struct {
	Email string `json:"emailClient"`
	Key   []byte `json:"key"`
}

type dataMesage struct {
	RequestType byte   `json:"requestType"`
	Message     []byte `json:"message"`
}

func (c *Client) readPump() {
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
			log.Println(err)
			if c.logginData == nil {
				c.registr <- 100
			}
			break
		}

		if c.logginData == nil {
			err = c.logickLoggin(datMessage)
			if err != nil {
				c.registr <- 100
				log.Println(err)
				break
			} else {
				// Авторизация прошла успешно
				c.registr <- 100

				log.Println("AUTARIZATION")
				parseNewClient, _ := json.Marshal(c.logginData)

				var requstMessage dataMesage
				requstMessage.RequestType = 101
				requstMessage.Message = parseNewClient

				sendMessage, _ := json.Marshal(requstMessage)

				c.send <- []byte(string(sendMessage))
			}
		}
	}
}

func (c *Client) writePump() {
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
func serveWs(hub *Hub, w http.ResponseWriter, r *http.Request) {

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	client := &Client{hub: hub, conn: conn, send: make(chan []byte, 256), logginData: nil}
	client.hub.register <- client

	go client.writePump()
	go client.readPump()
	go client.waitAutentification()
}

func (c *Client) waitAutentification() {
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

func (c *Client) disconnect() {
	c.hub.unregister <- c
	c.conn.Close()
}

func typeMessage(message *[]byte) (*dataMesage, error) {
	var inputMessage dataMesage

	err := json.Unmarshal(*message, &inputMessage)
	if err != nil {
		return &inputMessage, err
	}

	return &inputMessage, nil
}

func (c *Client) logickLoggin(datMessage *dataMesage) error {
	var retunrError error

	// Достаем данные из сообщения
	err := json.Unmarshal(datMessage.Message, &c.logginData)
	if err != nil {
		retunrError = err
		return err
	}

	if c.logginData == nil {
		retunrError = errors.New("Error logginData 101")
		return retunrError
	}

	switch datMessage.RequestType {
	case 100: // Запрос проверки логина пароля

		// Проверяем наличие не нулевых данных в сообщении
		if c.logginData.Email == "" || len(c.logginData.Key) == 0 {
			retunrError = errors.New("error email or key")
			break
		} else {
			// Идем в БД для проверки данных
			bdInfo, err := checkUser(*c.logginData)
			if err != nil {
				retunrError = err
				break
			}
			// Если бд возвращает пустого клиента, говорим что логин и пароль не верны
			if bdInfo == nil {
				retunrError = errors.New("dont accaunt")
				break
			} else {
				// Если клиент с таким логином есть, проверяем пароль

				hash, err := hashSum(c.logginData.Key, bdInfo.salt)
				if err != nil {
					retunrError = err
					break
				}

				log.Println(c.logginData.Key)
				log.Println(hash)
				log.Println(bdInfo.key)

				// Проверям ключи, если они не верны говорим что не верный пароль
				if bytes.Equal(hash, bdInfo.key) == false {
					retunrError = errors.New("error key")
					break
				}
				log.Println("Connect")
			}
			break
		}
		break

	case 101: // Запрос регистрации
		// Временная реализация для проверки коннекта

		// Проверяем что бы логин не был нулевым
		if c.logginData.Email == "" {
			retunrError = errors.New("null email")
			break
		} else {
			// Создаем нового клиента и добавляем его в БД

			content, err := getUserInfo(c.logginData.Email)
			if err != nil {
				retunrError = err
				break
			}

			var infoUser contentParse

			err = json.Unmarshal(content, &infoUser)
			if err != nil {
				retunrError = err
				break
			}

			if infoUser.Email == "" {
				retunrError = errors.New("nil email content")
				break
			} else {
				c.logginData.Email = infoUser.Email
			}

			//Проверяем есть ли в БД пользователь с таким логином
			bdInfo0, err := checkUser(*c.logginData)
			if err != nil {
				retunrError = err
				break
			}

			// Генерируем ключ для него
			newKey, err := randGenerate(64)
			if err != nil {
				retunrError = err
				break
			}

			// Генерируем соль для ключа
			newSalt, err := randGenerate(8)
			if err != nil {
				retunrError = err
				break
			}

			c.logginData.Key = newKey

			hash, err := hashSum(newKey, newSalt)
			if err != nil {
				retunrError = err
				break
			}

			infoBd := &bdInfo{email: c.logginData.Email, key: hash, salt: newSalt}

			// Если пользователя нет, то регестрируем его
			if bdInfo0 == nil {
				// Добавляем в бд нового пользователь
				err = addUser(infoBd)
				if err != nil {
					retunrError = err
					break
				}

				log.Println("Regestration")

				break
			} else {
				// Обновляем в бд пользователя
				err = updateUser(infoBd)
				if err != nil {
					retunrError = err
					break
				}
				// Такой аккаунт уже зарегестрирован
				log.Println("Re Regestration")
				break
			}
		}
		//retunrError = errors.New("101")
		break

	default:
		retunrError = errors.New("default")
		break
	}
	return retunrError
}

func randGenerate(lenght byte) ([]byte, error) {
	b := make([]byte, lenght)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}
	return b, err
}

func hashSum(key []byte, salt []byte) ([]byte, error) {
	h := sha256.New()
	key0 := make([]byte, 0)
	key0 = append(key0, key...)
	key0 = append(key0, salt...)
	_, err := h.Write(key0)
	if err != nil {
		return nil, err
	}

	generateKey := h.Sum(nil)
	return generateKey, err
}

func getUserInfo(code string) ([]byte, error) {

	googleOauthConfig = &oauth2.Config{
		RedirectURL:  "http://localhost:5000/dummy/oauth2callback",
		ClientID:     "93334427286-9kkusop9sjl32iml5qasuc58dhht25q7.apps.googleusercontent.com", //os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: "UdoSX27W3kMgERgTY9lqH6Qv",                                                //os.Getenv("GOOGLE_CLIENT_SECRET"),
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
		Endpoint:     google.Endpoint,
	}

	token, err := googleOauthConfig.Exchange(oauth2.NoContext, code)
	if err != nil {
		return nil, err
	}

	response, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return contents, nil
}

type contentParse struct {
	Id             string `json:"id"`
	Email          string `json:"email"`
	Verified_email bool   `json:"verified_email"`
	Name           string `json:"name"`
	Given_name     string `json:"given_name"`
	Link           string `json:"link"`
	Picture        string `json:"picture"`
}
