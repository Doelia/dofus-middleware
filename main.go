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

	// GA90390069329;90069329
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


func encodeChar(ch int) uint8 {
	alphabet := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-_"
	return alphabet[ch]
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


func main() {
	go login()
	go game()
	StartWebSocket()
	//themap := getMap(710)
	//path := AStar(themap, 76, 433)
	//fmt.Println("path", path)
	//encoded := encodePath(themap, path)
	//fmt.Println("encoded", encoded)
	//fmt.Println("decoded", decodePath(encoded))
}