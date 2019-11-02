package main

import (
	"dofusmiddleware/database"
	"dofusmiddleware/world"
	"fmt"
)

func main() {
	//go StartRealmProxy()
	//go StartGameProxy()
	//web.StartWebSocket()

	world.MapWithTriggers = database.GetMapriggers()
	fmt.Println(world.MapWithTriggers)

	//themap := getMap(710)
	//path := AStar(themap, 76, 433)
	//fmt.Println("path", path)
	//encoded := encodePath(themap, path)
	//fmt.Println("encoded", encoded)
	//fmt.Println("decoded", decodePath(encoded))
}