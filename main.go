package main

import (
	"bytes"
	"fmt"
	"net"
	"os/exec"
	"strings"
)

var Characters []Character

type Character struct {
	Name string
	IdCharDofus string
	Id string
	ConnClient net.Conn
	ConnServer net.Conn
}

func isOneOfMyCharacter(name string) bool {
	for _, c := range Characters {
		if c.Name == name || c.IdCharDofus == name {
			return true
		}
	}
	return false
}

func getChararacter(id string) *Character {
	for i, c := range Characters {
		if c.Id == id {
			return &Characters[i]
		}
	}
	return nil
}

func login() {

	fmt.Print("Hello world\n")

	p := Server{
		Addr:   "127.0.0.1:478",
		Target: "34.251.172.139:443",
		ModifyResponse: func(b *[]byte, id string) {
			packet := string(*b)
			//fmt.Println("[login] server->client: " + packet)
			if strings.HasPrefix(packet, "AXK3413389?ag7") {
				token := packet[14:]
				fmt.Println("bibly packet transform!,  token=" + token)
				by := bytes.NewBufferString("AYK127.0.0.1:5555;" + token)
				*b = by.Bytes()
				//fmt.Println(*b)
			}
		},
	}

	err := p.ListenAndServe()
	if err != nil {
		fmt.Println(err)
	}

}

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

func processPerso(id string, packet string) {
	splited := strings.Split(packet, "|")
	pr := splited[2]

	params := strings.Split(pr, ";")
	name := params[1]

	fmt.Println("Personnage " + name)
	var char = getChararacter(id)
	char.Name = name
	char.IdCharDofus = params[0]
}

func startTurn(id string, packet string) {
	splited := strings.Split(packet[3:], "|")
	idCharTurn := splited[0]
	char := getChararacter(id)
	if char.IdCharDofus == idCharTurn {
		fmt.Println("Start turn of " + char.Name)

		cmd := "/Users/stephane/Documents/dev/perso/dofus/" + char.Name + ".sh"
		out := exec.Command("/bin/bash", cmd)
		_ = out.Run()
	}
}

// PIKDoelia|Lotahi
func popupInvitation(id string, packet string) {
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
func popupExchange(id string, packet string) {
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

func game() {

	fmt.Print("Hello world\n")

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
					fmt.Println("[" + char.Name + "] server->client: " + strPacket)
				}

				if strings.HasPrefix(string(p), "ALK") {
					processPerso(id, strPacket)
				}

				if strings.HasPrefix(string(p), "GTS") {
					startTurn(id, strPacket)
				}

				if strings.HasPrefix(string(p), "PIK") {
					popupInvitation(id, strPacket)
				}

				if strings.HasPrefix(string(p), "ERK") {
					popupExchange(id, strPacket)
				}
			}

		},
		ModifyRequest: func(b *[]byte, id string) {
			strPacket := string(*b)
			strPacket = strPacket[:len(strPacket) - 1] // Remove last '0' byte

			char := getChararacter(id)
			if char != nil {
				fmt.Println("[" + char.Name + "] client->server: " + strPacket)
			}

			//if (strings.HasPrefix(packet, ""))
		},
	}

	err := p.ListenAndServe()
	if err != nil {
		fmt.Println(err)
	}

}

func main() {
	go login()
	game()
}