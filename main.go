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


func main() {
	go StartWebSocket()
	go login()
	go game()
	InputKeyboard()
}