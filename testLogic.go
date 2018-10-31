package main

import (
	"log"
	"serviceConnection"
	"serviceInfoPlayer"
)

// Структура клиента
type Client struct {
	*serviceConnection.Сlient
	infoPlayer *serviceInfoPlayer.InfoPlayer
}

// Вызывается, если клиент прошел авторизацию
func (c *Client) AutorizationCompleted() {
	c.infoPlayer = &serviceInfoPlayer.InfoPlayer{Gold: -1}
	err := c.infoPlayer.LoadInfo(c.Id, db)
	if err != nil {
		log.Println(err)
		c.Disconnect()
	}
	data, err := c.infoPlayer.ReturnDataInfo()
	if err != nil {
		c.Disconnect()
	}
	c.Write(110, data)
	//log.Println(string(*data))
}

// Вызывается, если клиент отключился
func (c *Client) ClientDisconnect() {
	log.Println("Client Disconnect")
}

// Вызывается при инициализации структуры, до авторизации
func (c *Client) Inicialization(connectionClient *serviceConnection.Сlient) {
	c.Сlient = connectionClient
}

// Принимает входящие пакеты
func (c *Client) Read(typeRequest byte, data []byte) {

}

// 100 - Прием: Подключение
// 101 - Прием: Регистрация <---> Отправка: Информация о подключенном клиенте
// 110 - Отправка: Информация о ресурсах игрока
