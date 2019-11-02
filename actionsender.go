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

func joinFightCharacter(char Character, startedBy string) {
	time.Sleep(time.Duration(500) * time.Millisecond)

	// GA90390069329;90069329
	packetConfirm := "GA903" + startedBy + ";" + startedBy
	fmt.Println("send join fight packet to " + char.Name)

	sendPacket(char, packetConfirm)
}

func readyFightCharacter(char Character) {
	time.Sleep(time.Duration(1000) * time.Millisecond)

	packetConfirm := "GR1"
	fmt.Println("send ready fight packet to " + char.Name)

	sendPacket(char, packetConfirm)
}

func sendPacket(char Character, packet string) {
	fmt.Println("send packet", packet, "to", char.Name)
	packetConfirm := bytes.NewBufferString(packet)
	packetConfirm.WriteByte(0)
	packetConfirm.WriteString("\n")
	_, _ = char.ConnServer.Write(packetConfirm.Bytes())
}
