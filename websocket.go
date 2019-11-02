package main

import (
	"encoding/json"
	"fmt"
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


func SendOptions(options OptionsStruct) {
	message, _ := json.Marshal(options)
	SendPacket("OPTIONS", string(message))
}

func SendMap() {
	mapp := buildCellsFromMapData("HhaaeaaaaaHhaaeaaaaaHhaaeaaaaaHhaaeaaadyHhaaeaaaaaHhaaeaaadyHhaaeaaadyHhGaeaaaaaHhaaeaaaaaHhaaeaaaaaHhaaeaaaaaHhaaeaaadyHhaaeaaaaaHhaaeaaaaaHhaaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhG9eaaaaaHhGaeaaaaaHhaaeaaadyHhaaeaaadyHhGaeaaaaaHhqaeqgaaaHhGaeaaaaaHhG9eaaaaaHhGaeaaaaaHhqaeaaaaaHhaaeaaaaaGha9eaaaaRHhGaeaaaaaHhG9eaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhG9eaaaaaHhG9eaaaaaHhGaeaaaaaHhGaeaaaeBHhGaeaaaaaHhGaeaaaaaHha9eaaadyHhaaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHha9eaaaaOHhGaeaaaeBHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHha9eaaaaRHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhaaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhaaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhaaeaaaaaHhGaeaaaaaHhGaeiFaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhaaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhaaeaaadyHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhaaeaaaaaHhGaeaaaaaHhGaeaaaeAHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeiFaaaHhGaeaaaaaHhGaeaaaaaHhGaeiFaaaHhGaeaaaaaHhGaeaaaaaHhaaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaem0aaaHha9eaaaiSHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhaaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhG9eaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaGha9eaaadzHhGaeaaaaaHhGaeb4aaaHhGaeb4aaaHhaaeb4aaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaem0aaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaem0aaaHhGaeaaaaaHhGaeb5WaaHhGkeaaaaaHhGkeqgaaaHhaaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeiFaaaHhGaem0aaaHhGaemHWaaHhGkeaaaaaHhGkeaaaaaHhakeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeiFaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGkeaaaaaHhGkeaaaaaHhGaeb4GaaHhaaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeb4WaaHhGkeaaaaaHhGaemHqaaHhaaemHGaaHhGaeqgaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeiFaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaemHWaaHhGkeaGaaaHhGaemHqaaHhGaeaaaaaHhaaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeiFaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGkeaaaaaHhGaeb4qaaHhGaeaaaaaHhaaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeb4WaaHhGkeaaaaaHhGaeaaaaaHhGaeaaaaaHhaaeaaaaaHhGaeiFaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaemHWaaHhGkeaaaaaHhGaemHqaaHhGaeaaaaaHha7eaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaem0aaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGkeaaaaaHhGaemHqaaHhGaeaaaaaHhHfeaaaaaHhaaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeb4WaaHhGaeb5qaaHhGaeaaaaaHhGaeaaaaaHhbfeaaaiGHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaem0aaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaemHWaaHhGkeaaaaaHhGaeaaaaaHhG5eaaaaaHhHfeaaaaaHhaaeaaaaaHhGaeaaadmHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaem0aaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGkeaaaaaHhGaemHqaaHhG5eaaaaaHhHfeaaaaaHhaaeaaaaaHxG5eaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeb5WaaHhGaeb4qaaH3G7eaaaaaHhHfeaaaaaHhGaeaaaaaHhaaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaemHWaaHhGkeaaaaaHhGaeaaaaaHhG6eaaadmGhaaeaaaaOHhaaeaaaaaHhHfeaaaaaHhGaeaaaaaHhHfeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeb5aaaHhGkeaaaaaHhGaemHqaaHhGaeaaaaaHhHfeaaaaaHhGaeaaaaaHhbfeaaadtHhGaeaaaaaHhHfeaaaaaHhGaeaaaeBHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeb5WaaHhGkeaaaaaHhGaemHqaaHhGaeaaaaaHhHfeaaaaaHhGaeaaaaaHhaaeaaaaaHhbfeaaaiSHhGaeaaaaaHhHfeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeiFaaaHhGaeaaaaaHhGaemHWaaHhGkeaaaaaHhGaeb4qaaHhGaeaaaaaHhHfeaaaaaHhHfeaaaaaHhGNeaaaaaHhbfeaaaaaHhHfeaaaaaGhbfeaaadxHhaaeaaaaaHhaaeaaaiGHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGkeaaaaaHhGkeaaaaaHhGaeaaaaaHha5eaaaiGHhHfeaaaaaGhbfeaaaaRHhbfeaaaaOHhGNeaaaaaGhbfeaaaaRHhbfeaaaiMHhaaeaaaaaHxG5eaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeb5WaaHhGkeqgaaaHhGaeb4qaaHhG5eaaaaaHhbfeaaaaSHhGNeaaaaaHhHfeaaaaaHhbfeaaaaaHhbfeaaaaaHhaNeaaaaaHhbfeaaadtHhbfeaaaaaHxa7eaaaaaHhaaeaaaaaHhaaeaaaaaHhakeaaaaaHhakeaaaaaH3a7eaaaaaHhaNeaaaaaHhbfeaaaaaHhbfeaaaaaHhbfeaaaaa")
	message, _ := json.Marshal(mapp)
	SendPacket("MAP", string(message))

	message3, _ := json.Marshal(Path)
	SendPacket("PATH", string(message3))
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
}

func SendPacket(typepacket string, content string) {
	if GuiSocket != nil {
		err := GuiSocket.WriteMessage(websocket.TextMessage, []byte(typepacket + "|" + content))
		if err != nil {
			log.Println("write error: ", err)
		}
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
	SendOptions(Options)
	SendMap()
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