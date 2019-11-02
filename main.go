package main

import (
	"bytes"
	"fmt"
	"time"
)


func extractPackets(b* []byte) [][]byte {
	var packets [][]byte
	current := bytes.NewBuffer([]byte{})

	by := bytes.NewBuffer(*b)

	for {
		bi, err := by.ReadByte()
		if err != nil {
			return packets
		}
		current.WriteByte(bi)
		if bi == 0 {
			packets = append(packets, current.Bytes())
			current = bytes.NewBuffer([]byte{})
		}
	}
}


func joinFightCharacter(char Character, startedBy string) {
	time.Sleep(time.Duration(500) * time.Millisecond)

	//GA90390069329;90069329
	packetConfirm := bytes.NewBufferString("GA903" + startedBy + ";" + startedBy)
	fmt.Println("send join fight packet to " + char.Name)

	packetConfirm.WriteByte(0)
	packetConfirm.WriteString("\n")
	_, _ = char.ConnServer.Write(packetConfirm.Bytes())
}

func readyFightCharacter(char Character) {
	time.Sleep(time.Duration(1000) * time.Millisecond)

	packetConfirm := bytes.NewBufferString("GR1")
	fmt.Println("send ready fight packet to " + char.Name)

	packetConfirm.WriteByte(0)
	packetConfirm.WriteString("\n")
	_, _ = char.ConnServer.Write(packetConfirm.Bytes())
}

func moveChar(char Character, packet string, counter int) {
	fmt.Println("send move " + packet + " to " + char.Name)
	time.Sleep(time.Duration(counter * 200) * time.Millisecond)
	packetConfirm := bytes.NewBufferString(packet)
	packetConfirm.WriteByte(0)
	packetConfirm.WriteString("\n")
	_, _ = char.ConnServer.Write(packetConfirm.Bytes())
}

func decodeChar(c uint8) uint8 {
	if c >= 'a' && c <= 'z' {
		return c - 'a'
	}
	if c >= 'A' && c <= 'Z' {
		return (c - 'A') + 26
	}
	if c >= '0' && c <= '9' {
		return (c - '0') + 26 * 2
	}
	if c == '-' {
		return 62
	}
	if c == '_' {
		return 63
	}
	return 0
}

func decodePath(path string) uint8 {
	lastCell := path[len(path)-2:]
	c1 := lastCell[0]
	c2 := lastCell[1]
	return decodeChar(c1) * 64 + decodeChar(c2)
}

var Path []int

func main() {
	//go login()
	//go game()
	//InputKeyboard()

	cells := buildCellsFromMapData("HhaaeaaaaaHhaaeaaaaaHhaaeaaaaaHhaaeaaadyHhaaeaaaaaHhaaeaaadyHhaaeaaadyHhGaeaaaaaHhaaeaaaaaHhaaeaaaaaHhaaeaaaaaHhaaeaaadyHhaaeaaaaaHhaaeaaaaaHhaaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhG9eaaaaaHhGaeaaaaaHhaaeaaadyHhaaeaaadyHhGaeaaaaaHhqaeqgaaaHhGaeaaaaaHhG9eaaaaaHhGaeaaaaaHhqaeaaaaaHhaaeaaaaaGha9eaaaaRHhGaeaaaaaHhG9eaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhG9eaaaaaHhG9eaaaaaHhGaeaaaaaHhGaeaaaeBHhGaeaaaaaHhGaeaaaaaHha9eaaadyHhaaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHha9eaaaaOHhGaeaaaeBHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHha9eaaaaRHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhaaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhaaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhaaeaaaaaHhGaeaaaaaHhGaeiFaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhaaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhaaeaaadyHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhaaeaaaaaHhGaeaaaaaHhGaeaaaeAHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeiFaaaHhGaeaaaaaHhGaeaaaaaHhGaeiFaaaHhGaeaaaaaHhGaeaaaaaHhaaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaem0aaaHha9eaaaiSHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhaaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhG9eaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaGha9eaaadzHhGaeaaaaaHhGaeb4aaaHhGaeb4aaaHhaaeb4aaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaem0aaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaem0aaaHhGaeaaaaaHhGaeb5WaaHhGkeaaaaaHhGkeqgaaaHhaaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeiFaaaHhGaem0aaaHhGaemHWaaHhGkeaaaaaHhGkeaaaaaHhakeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeiFaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGkeaaaaaHhGkeaaaaaHhGaeb4GaaHhaaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeb4WaaHhGkeaaaaaHhGaemHqaaHhaaemHGaaHhGaeqgaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeiFaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaemHWaaHhGkeaGaaaHhGaemHqaaHhGaeaaaaaHhaaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeiFaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGkeaaaaaHhGaeb4qaaHhGaeaaaaaHhaaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeb4WaaHhGkeaaaaaHhGaeaaaaaHhGaeaaaaaHhaaeaaaaaHhGaeiFaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaemHWaaHhGkeaaaaaHhGaemHqaaHhGaeaaaaaHha7eaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaem0aaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGkeaaaaaHhGaemHqaaHhGaeaaaaaHhHfeaaaaaHhaaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeb4WaaHhGaeb5qaaHhGaeaaaaaHhGaeaaaaaHhbfeaaaiGHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaem0aaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaemHWaaHhGkeaaaaaHhGaeaaaaaHhG5eaaaaaHhHfeaaaaaHhaaeaaaaaHhGaeaaadmHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaem0aaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGkeaaaaaHhGaemHqaaHhG5eaaaaaHhHfeaaaaaHhaaeaaaaaHxG5eaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeb5WaaHhGaeb4qaaH3G7eaaaaaHhHfeaaaaaHhGaeaaaaaHhaaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaemHWaaHhGkeaaaaaHhGaeaaaaaHhG6eaaadmGhaaeaaaaOHhaaeaaaaaHhHfeaaaaaHhGaeaaaaaHhHfeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeb5aaaHhGkeaaaaaHhGaemHqaaHhGaeaaaaaHhHfeaaaaaHhGaeaaaaaHhbfeaaadtHhGaeaaaaaHhHfeaaaaaHhGaeaaaeBHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeb5WaaHhGkeaaaaaHhGaemHqaaHhGaeaaaaaHhHfeaaaaaHhGaeaaaaaHhaaeaaaaaHhbfeaaaiSHhGaeaaaaaHhHfeaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeiFaaaHhGaeaaaaaHhGaemHWaaHhGkeaaaaaHhGaeb4qaaHhGaeaaaaaHhHfeaaaaaHhHfeaaaaaHhGNeaaaaaHhbfeaaaaaHhHfeaaaaaGhbfeaaadxHhaaeaaaaaHhaaeaaaiGHhGaeaaaaaHhGaeaaaaaHhGaeaaaaaHhGkeaaaaaHhGkeaaaaaHhGaeaaaaaHha5eaaaiGHhHfeaaaaaGhbfeaaaaRHhbfeaaaaOHhGNeaaaaaGhbfeaaaaRHhbfeaaaiMHhaaeaaaaaHxG5eaaaaaHhGaeaaaaaHhGaeaaaaaHhGaeb5WaaHhGkeqgaaaHhGaeb4qaaHhG5eaaaaaHhbfeaaaaSHhGNeaaaaaHhHfeaaaaaHhbfeaaaaaHhbfeaaaaaHhaNeaaaaaHhbfeaaadtHhbfeaaaaaHxa7eaaaaaHhaaeaaaaaHhaaeaaaaaHhakeaaaaaHhakeaaaaaH3a7eaaaaaHhaNeaaaaaHhbfeaaaaaHhbfeaaaaaHhbfeaaaaa")
	mapp := Map{
		Width: 15,
		Height:17,
		Cells: cells,
	}

	Path = AStar(mapp, 170,450)

	StartWebSocket()

	//voisions := getNeighbors(mapp, getCellOfMap(mapp, 260))
	//fmt.Println(voisions)
}