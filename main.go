package main

import (
	"dofusmiddleware/database"
	"dofusmiddleware/websocket"
	"dofusmiddleware/world"
)

var web websocket.WebSocket

func main() {
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