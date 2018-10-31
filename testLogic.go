package main

import (
	"log"
	"serviceConnection"
	"serviceInfoPlayer"
)

type Client struct {
	*serviceConnection.Сlient
	infoPlayer *serviceInfoPlayer.InfoPlayer
}

func (c *Client) AutorizationCompleted() {
	c.infoPlayer = &serviceInfoPlayer.InfoPlayer{Gold: -1}
	err := c.infoPlayer.LoadInfo(c.Id, db)
	if err != nil {
		log.Println(err)
		c.Disconnect()
	}
	//	err = c.infoPlayer.SaveInfo(c.Id, db)
	//	if err != nil {
	//		c.Disconnect()
	//	}
	log.Println(c.infoPlayer)
	log.Println(c.infoPlayer.ReturnDataInfo())
}

func (c *Client) ClientDisconnect() {
	log.Println("Client Disconnect")
}

func (c *Client) Inicialization(connectionClient *serviceConnection.Сlient) {
	c.Сlient = connectionClient
}

func (c *Client) Read(typeRequest byte, data []byte) {

}
