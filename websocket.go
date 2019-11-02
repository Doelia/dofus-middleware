package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"strconv"
	"strings"
)

var upgrader = websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true} } // use default options
var GuiSocket *websocket.Conn

func SendCharacters(characters []Character) {
	message, _ := json.Marshal(characters)
	SendPacket("CHARACTERS", string(message))
}


func SendOptions(options OptionsStruct) {
	message, _ := json.Marshal(options)
	SendPacket("OPTIONS", string(message))
}

func SendMap(themap Map) {
	message, _ := json.Marshal(themap)
	SendPacket("MAP", string(message))
}

func SendPath(path []int) {
	message, _ := json.Marshal(path)
	SendPacket("PATH", string(message))
}


func onMessage(packet string) {
	parts := strings.Split(packet, "|")
	typepacket := parts[0]
	if typepacket == "FOCUS" {
		go SwitchToCharacter(parts[1])
	}
	if typepacket == "SET_OPTION" {
		fmt.Println("websocket input : SET_OPTIONS " + parts[1] + " " + parts[2])
		if parts[1] == "ShowInputPackets" {
			Options.ShowInputPackets = parts[2] == "true"
		}
		if parts[1] == "ShowOutputPackets" {
			Options.ShowOutputPackets = parts[2] == "true"
		}
		if parts[1] == "DispatchMoves" {
			Options.DispatchMoves = parts[2] == "true"
		}
		if parts[1] == "FocusWindowOnCharacterTurn" {
			Options.FocusWindowOnCharacterTurn = parts[2] == "true"
		}
		if parts[1] == "AutoJoinFight" {
			Options.AutoJoinFight = parts[2] == "true"
		}
		if parts[1] == "AutoReadyFight" {
			Options.AutoReadyFight = parts[2] == "true"
		}
		SendOptions(Options)
	}
	if typepacket == "SET_CHARACTER_OPTION" {
		fmt.Println("websocket input : SET_CHARACTER_OPTIONS " + parts[1] + " " + parts[2] + " " + parts[3])
		if parts[2] == "OptionAutoPassTurn" {
			getChararacter(parts[1]).OptionAutoPassTurn = parts[3] == "true"
		}
		SendCharacters(Characters)
	}
	if typepacket == "PROCESS_PATH" {
		idMap, _ := strconv.Atoi(parts[1])
		cellStart, _ := strconv.Atoi(parts[2])
		cellEnd, _ := strconv.Atoi(parts[3])

		fmt.Println("process path", idMap, cellStart, cellEnd)

		themap := getMap(idMap)
		path := AStar(themap, cellStart, cellEnd)
		encodedPath := encodePath(themap, path)

		SendPath(path)
		sendMovePacket(*getAConnectedCharacter(), encodedPath)
	}
}

func SendPacket(typepacket string, content string) {
	if GuiSocket != nil {
		err := GuiSocket.WriteMessage(websocket.TextMessage, []byte(typepacket + "|" + content))
		if err != nil {
			log.Println("write error: ", err)
		}
	}
}

func OnConnect() {
	SendCharacters(Characters)
	SendOptions(Options)
	themap := getMap(getAConnectedCharacter().MapId)
	SendMap(themap)
}

func echo(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	GuiSocket = c
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()
	OnConnect()
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