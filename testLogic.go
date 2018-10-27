package main

import (
	"encoding/json"
	"log"
	"serviceAutorization"
	"serviceConnection"
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
		err := serviceAutorization.Autorization(datMessage.RequestType, datMessage.Message, db)
		if err != nil {
			log.Println(err)
			c.connectionClient.Disconnect()
			return
		} else {
			c.autorization = true
		}
	} else {

	}
}

func (c *Client) Write(data []byte) {

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
