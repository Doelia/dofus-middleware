package main

import (
	"bytes"
	"fmt"
	"strings"
)

func OnCharacterEnterInGame(id string, packet string) {
	splited := strings.Split(packet, "|")
	pr := splited[2]

	params := strings.Split(pr, ";")
	name := params[1]

	fmt.Println("Character enter in game : " + name)
	var char = getChararacter(id)
	char.Name = name
	char.IdCharDofus = params[0]
}

func OnStartTurn(id string, packet string) {
	splited := strings.Split(packet[3:], "|")
	idCharTurn := splited[0]
	char := getChararacter(id)
	if char.IdCharDofus == idCharTurn {
		fmt.Println("Start turn of " + char.Name)
		SwitchToCharacter(char.Name)
	}
}

// PIKDoelia|Lotahi
func OnPopupGroupInvitation(id string, packet string) {
	splited := strings.Split(packet[3:], "|")
	inviter := splited[0]
	invited := splited[1]

	char := getChararacter(id)
	fmt.Println(inviter + " " + invited + " " + char.Name)

	// Im invited
	if invited == char.Name {
		if isOneOfMyCharacter(inviter) {
			fmt.Println("Im ("+ invited +") invited to join "+ inviter +" group's")
			packetConfirm := bytes.NewBufferString("PA")
			packetConfirm.WriteByte(0)
			packetConfirm.WriteString("\n")
			_, _ = char.ConnServer.Write(packetConfirm.Bytes())
		}
	}
}

//  ERK90069329|90069284|1
func OnPopupExchange(id string, packet string) {
	splited := strings.Split(packet[3:], "|")
	inviter := splited[0]
	invited := splited[1]

	char := getChararacter(id)
	fmt.Println(inviter + " " + invited + " " + char.Name)

	// Im invited
	if invited == char.IdCharDofus {
		if isOneOfMyCharacter(inviter) {
			fmt.Println("Im ("+ invited +") invited to exchange with "+ inviter)
			packetConfirm := bytes.NewBufferString("EA")
			packetConfirm.WriteByte(0)
			packetConfirm.WriteString("\n")
			_, _ = char.ConnServer.Write(packetConfirm.Bytes())
		}
	}
}

func OnMoveCharater(id string, packet string) {
	//counter := 0
	//for _, c := range Characters {
	//	if c.Name != "" && c.Id != id {
	//		counter := counter + 1
	//		go moveChar(c, packet, counter)
	//	}
	//}
}

func game() {

	fmt.Print("New Game proxy")

	p := Server{
		Addr:   "127.0.0.1:5555",
		Target: "52.19.56.159:443",
		ModifyResponse: func(b *[]byte, id string) {
			//fmt.Println(*b)

			packets := extractPackets(b)
			for _, p := range packets {

				strPacket := string(p)
				strPacket = strPacket[:len(strPacket) - 1] // Remove trailing '0' byte

				char := getChararacter(id)
				if char != nil {
					fmt.Println("[" + char.Name + "] WebSocket->client: " + strPacket)
				}

				if strings.HasPrefix(string(p), "ALK") {
					OnCharacterEnterInGame(id, strPacket)
				}

				if strings.HasPrefix(string(p), "GTS") {
					OnStartTurn(id, strPacket)
				}

				if strings.HasPrefix(string(p), "PIK") {
					OnPopupGroupInvitation(id, strPacket)
				}

				if strings.HasPrefix(string(p), "ERK") {
					OnPopupExchange(id, strPacket)
				}
			}

		},
		ModifyRequest: func(b *[]byte, id string) {

			bytess := make([]byte, len(*b))
			copy(bytess, *b)
			bytess = bytess[:len(bytess) - 1] // Remove trailing '\n' byte
			bytess[len(bytess) - 1] = 0

			packets := extractPackets(&bytess)
			for _, p := range packets {

				strPacket := string(p)
				strPacket = strPacket[:len(strPacket) - 1] // Remove trailing '0' byte

				char := getChararacter(id)
				if char != nil {
					fmt.Println("[" + char.Name + "] client->WebSocket: " + strPacket)
				}

				if strings.HasPrefix(strPacket, "GA001") {
					OnMoveCharater(id, strPacket)
				}
			}

		},
	}

	err := p.ListenAndServe()
	if err != nil {
		fmt.Println(err)
	}

}
