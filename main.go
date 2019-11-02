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

	web = websocket.WebSocket{
		OnConnexion: OnSocketConnexion,
		OnMessage: OnSocketMessage,
	}
	web.StartWebSocket()

	//themap := getMap(710)
	//path := AStar(themap, 76, 433)
	//fmt.Println("path", path)
	//encoded := encodePath(themap, path)
	//fmt.Println("encoded", encoded)
	//fmt.Println("decoded", decodePath(encoded))
}