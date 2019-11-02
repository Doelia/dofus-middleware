package socket

import (
	"bytes"
	"dofusmiddleware/world"
	"fmt"
	"time"
)

func MoveChar(connexion world.Connexion, packet string, counter int) {
	fmt.Println("send move " + packet + " to " + connexion.Player.Name)
	time.Sleep(time.Duration(counter * 200) * time.Millisecond)
	sendPacket(connexion, packet)
}

func SendMovePacket(connexion world.Connexion, path string) {
	sendPacket(connexion, "GA001" + path)
}

func JoinFightCharacter(connexion world.Connexion, startedBy string) {
	time.Sleep(time.Duration(500) * time.Millisecond)

	// GA90390069329;90069329
	packetConfirm := "GA903" + startedBy + ";" + startedBy
	fmt.Println("send join fight packet to " + connexion.Player.Name)

	sendPacket(connexion, packetConfirm)
}

func ReadyFightCharacter(connexion world.Connexion) {
	time.Sleep(time.Duration(1000) * time.Millisecond)

	packetConfirm := "GR1"
	fmt.Println("send ready fight packet to " + connexion.Player.Name)

	sendPacket(connexion, packetConfirm)
}

func SendConfirmAction(connexion world.Connexion) {
	packetConfirm := "PA"
	sendPacket(connexion, packetConfirm)
}

func SendPassTurn(connexion world.Connexion) {
	sendPacket(connexion, "GT")
}

func sendPacket(connexion world.Connexion, packet string) {
	fmt.Println("send packet", packet, "to", connexion.Player.Name)
	packetConfirm := bytes.NewBufferString(packet)
	packetConfirm.WriteByte(0)
	packetConfirm.WriteString("\n")
	_, _ = connexion.ConnServer.Write(packetConfirm.Bytes())
}

