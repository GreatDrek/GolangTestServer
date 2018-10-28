package main

import (
	"encoding/json"
	"log"
	"serviceAutorization"
	"serviceConnection"
	"time"
)

type Client struct {
	*serviceConnection.Сlient
	autorization bool
}

func (c *Client) Inicialization(connectionClient *serviceConnection.Сlient) {
	c.Сlient = connectionClient
}

func (c *Client) Read(data []byte) {
	datMessage, err := typeMessage(&data)
	if err != nil {
		c.Disconnect()
		log.Println(err)
		return
	}
	if c.autorization == false {
		log.Println("AutoStart", time.Now().String())
		logginData, err := serviceAutorization.Autorization(datMessage.RequestType, datMessage.Message, db)
		if err != nil {
			log.Println(err)
			c.Disconnect()
			return
		} else {
			parseNewClient, err := json.Marshal(logginData)
			if err != nil {
				c.Disconnect()
				return
			}
			c.WriteData(101, parseNewClient)
			c.autorization = true
			log.Println("AutoStop", time.Now().String())
		}
	} else {

	}
}

func (c *Client) WriteData(typeM byte, data []byte) {
	var requstMessage dataMesage
	requstMessage.RequestType = typeM
	requstMessage.Message = data

	sendMessage, err := json.Marshal(requstMessage)
	if err != nil {
		c.Disconnect()
		return
	}
	c.Write(sendMessage)
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
