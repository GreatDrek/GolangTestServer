package main

import (
	"encoding/json"
	"log"
	"time"
)

//type WorldOptions struct {
//	sizeX int
//	sizeY int
//}

type World struct {
	players map[*Player]bool

	register chan *Player

	unregister chan *Player
}

func newWorld() *World {
	return &World{
		register:   make(chan *Player),
		unregister: make(chan *Player),
		players:    make(map[*Player]bool),
	}
}

func (w *World) run() {
	for {
		select {
		case player := <-w.register:
			w.players[player] = true

			log.Println("New Player")

		case player := <-w.unregister:
			if _, ok := w.players[player]; ok {
				delete(w.players, player)
				//close(player.send)
			}

			log.Println("Delete Player")
		}
	}
}

func (w *World) BrodcastClients() {
	for {
		for player := range w.players {
			player.client.WriteData(201, []byte("Cheto"))
			//			!!!!!!!!
		}
		time.Sleep(time.Millisecond * 60)
	}
}

type Player struct {
	id     int
	posX   int
	posY   int
	client *Client
}

func (p *Player) Move(vecMove *Vector2D) {
	//!!!!!!!!!!!!!!!!!!!!!!!!!!
}

func (p *Player) InputPackage(typeMessage byte, data []byte) {
	if typeMessage == 200 {
		var vecMove Vector2D
		err := json.Unmarshal(data, &vecMove)
		if err != nil {
			p.client.Disconnect()
		} else {
			p.Move(&vecMove)
		}
	}
}

type Vector2D struct {
	X int `json:"x"`
	Y int `json:"y"`
}
