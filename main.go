package main

import "dofusmiddleware/web"

func main() {
	go StartRealmProxy()
	go StartGameProxy()
	web.StartWebSocket()
	//themap := getMap(710)
	//path := AStar(themap, 76, 433)
	//fmt.Println("path", path)
	//encoded := encodePath(themap, path)
	//fmt.Println("encoded", encoded)
	//fmt.Println("decoded", decodePath(encoded))
}