package main

import (
	"bytes"
	"fmt"
	"time"
)

func moveChar(char Character, packet string, counter int) {
	fmt.Println("send move " + packet + " to " + char.Name)
	time.Sleep(time.Duration(counter * 200) * time.Millisecond)
	sendPacket(char, packet)
}

func sendMovePacket(char Character, path string) {
	sendPacket(char, "GA001" + path)
}

func sendPacket(char Character, packet string) {
	fmt.Println("send packet", packet, "to", char.Name)
	packetConfirm := bytes.NewBufferString(packet)
	packetConfirm.WriteByte(0)
	packetConfirm.WriteString("\n")
	_, _ = char.ConnServer.Write(packetConfirm.Bytes())
}
