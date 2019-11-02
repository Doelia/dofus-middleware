package socket

import (
	"bytes"
	"dofusmiddleware/world"
	"fmt"
	"time"
)

func MoveChar(char world.Character, packet string, counter int) {
	fmt.Println("send move " + packet + " to " + char.Name)
	time.Sleep(time.Duration(counter * 200) * time.Millisecond)
	sendPacket(char, packet)
}

func SendMovePacket(char world.Character, path string) {
	sendPacket(char, "GA001" + path)
}

func JoinFightCharacter(char world.Character, startedBy string) {
	time.Sleep(time.Duration(500) * time.Millisecond)

	// GA90390069329;90069329
	packetConfirm := "GA903" + startedBy + ";" + startedBy
	fmt.Println("send join fight packet to " + char.Name)

	sendPacket(char, packetConfirm)
}

func ReadyFightCharacter(char world.Character) {
	time.Sleep(time.Duration(1000) * time.Millisecond)

	packetConfirm := "GR1"
	fmt.Println("send ready fight packet to " + char.Name)

	sendPacket(char, packetConfirm)
}

func sendPacket(char world.Character, packet string) {
	fmt.Println("send packet", packet, "to", char.Name)
	packetConfirm := bytes.NewBufferString(packet)
	packetConfirm.WriteByte(0)
	packetConfirm.WriteString("\n")
	_, _ = char.ConnServer.Write(packetConfirm.Bytes())
}
