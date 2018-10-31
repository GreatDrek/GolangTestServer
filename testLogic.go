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
	data, err := c.infoPlayer.ReturnDataInfo()
	if err != nil {
		c.Disconnect()
	}
	c.Write(110, data)
	log.Println(string(*data))
}

func (c *Client) ClientDisconnect() {
	log.Println("Client Disconnect")
}

func (c *Client) Inicialization(connectionClient *serviceConnection.Сlient) {
	c.Сlient = connectionClient
}

func (c *Client) Read(typeRequest byte, data []byte) {

}
