package dofusmiddleware

func main() {
	go Login()
	go Game()
	StartWebSocket()
	//themap := getMap(710)
	//path := AStar(themap, 76, 433)
	//fmt.Println("path", path)
	//encoded := encodePath(themap, path)
	//fmt.Println("encoded", encoded)
	//fmt.Println("decoded", decodePath(encoded))
}