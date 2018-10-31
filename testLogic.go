package main

import (
	"log"
	"serviceConnection"
)

type Client struct {
	*serviceConnection.Сlient
	infoPlayer *InfoPlayer
}

func (c *Client) AutorizationCompleted() {
	err := c.infoPlayer.LoadInfo(c.Id, db)
	if err != nil {
		log.Println(err)
		c.Disconnect()
	}
	log.Println(c.infoPlayer)
	c.infoPlayer.Gold = 100
	err = c.infoPlayer.SaveInfo(c.Id, db)
	if err != nil {
		c.Disconnect()
	}
	log.Println(c.infoPlayer)
}

func (c *Client) ClientDisconnect() {
	log.Println("Client Disconnect")
}

func (c *Client) Inicialization(connectionClient *serviceConnection.Сlient) {
	c.Сlient = connectionClient
}

func (c *Client) Read(typeRequest byte, data []byte) {

}
