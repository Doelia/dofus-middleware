package main

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"strings"
)

var upgrader = websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true} } // use default options
var GuiSocket *websocket.Conn

func SendCharacters(characters []Character) {
	message, _ := json.Marshal(characters)
	SendPacket("CHARACTERS", string(message))
}

func onMessage(packet string) {

	parts := strings.Split(packet, "|")
	typepacket := parts[0]
	content := parts[1]
	if typepacket == "FOCUS" {
		SwitchToCharacter(content)
	}
}

func SendPacket(typepacket string, content string) {
	err := GuiSocket.WriteMessage(websocket.TextMessage, []byte(typepacket + "|" + content))
	if err != nil {
		log.Println("write error: ", err)
	}
}

func echo(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	GuiSocket = c
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	SendCharacters(Characters)
	defer c.Close()
	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		if mt == websocket.TextMessage {
			onMessage(string(message))
		}
	}
}

func StartWebSocket() {
	http.HandleFunc("/ws", echo)
	log.Fatal(http.ListenAndServe("localhost:8001", nil))


}