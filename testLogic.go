package main

import (
	"log"
	"serviceConnection"
)

type Client struct {
	*serviceConnection.Сlient
}

func (c *Client) ClientDisconnect() {
	log.Println("Client Disconnect")
}

func (c *Client) Inicialization(connectionClient *serviceConnection.Сlient) {
	c.Сlient = connectionClient
}

func (c *Client) Read(typeRequest byte, data []byte) {

}
