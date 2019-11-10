package main

import (
	"dofusmiddleware/database"
	"dofusmiddleware/websocket"
	"dofusmiddleware/world"
	"math/rand"
	"time"
)

var web websocket.WebSocket

func main() {
	rand.Seed(time.Now().UnixNano())

	world.MapWithTriggers = database.GetMapriggers()

	go StartRealmProxy()
	go StartGameProxy()
	//
	web = websocket.WebSocket{
		OnConnexion: OnSocketConnexion,
		OnMessage: OnSocketMessage,
	}


	web.StartWebSocket()
	//encoded := encodePath(themap, path)
	//fmt.Println("encoded", encoded)
	//fmt.Println("decoded", decodePath(encoded))
}