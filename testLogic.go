package main

import (
	"encoding/json"
	"log"
	"serviceAutorization"
	"serviceConnection"
	"time"
)

type Client struct {
	connectionClient *serviceConnection.ConnectionClient
	autorization     bool
}

func (c *Client) Inicialization(connectionClient *serviceConnection.ConnectionClient) {
	c.connectionClient = connectionClient
}

func (c *Client) Read(data []byte) {
	datMessage, err := typeMessage(&data)
	if err != nil {
		c.connectionClient.Disconnect()
		log.Println(err)
		return
	}
	if c.autorization == false {
		logginData, err := serviceAutorization.Autorization(datMessage.RequestType, datMessage.Message, db)
		if err != nil {
			log.Println(err)
			c.connectionClient.Disconnect()
			return
		} else {
			parseNewClient, err := json.Marshal(logginData)
			if err != nil {
				c.connectionClient.Disconnect()
				return
			}
			c.Write(101, parseNewClient)
			c.autorization = true
			log.Println("Stop", time.Now().String())
		}
	} else {

	}
}

func (c *Client) Write(typeM byte, data []byte) {
	//parseNewClient, err := json.Marshal(c.logginData)

	var requstMessage dataMesage
	requstMessage.RequestType = typeM
	requstMessage.Message = data

	sendMessage, err := json.Marshal(requstMessage)
	if err != nil {
		c.connectionClient.Disconnect()
		return
	}
	c.connectionClient.Write(sendMessage)
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
