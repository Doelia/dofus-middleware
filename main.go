package main

import (
	"bufio"
	"bytes"
	"fmt"
	"net"
	"os"
	"os/exec"
	"strings"
	"time"
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

func getChararacter(search string) *Character {
	for i, c := range Characters {
		if c.Id == search || c.Name == search {
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

func moveChar(char Character, packet string, counter int) {
	fmt.Println("send move " + packet + " to " + char.Name)
	time.Sleep(time.Duration(counter * 200) * time.Millisecond)
	packetConfirm := bytes.NewBufferString(packet)
	packetConfirm.WriteByte(0)
	packetConfirm.WriteString("\n")
	_, _ = char.ConnServer.Write(packetConfirm.Bytes())
}

func outMoveCharater(id string, packet string) {
	counter := 0
	for _, c := range Characters {
		if c.Name != "" && c.Id != id {
			counter := counter + 1
			go moveChar(c, packet, counter)
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
					fmt.Println("[" + char.Name + "] client->server: " + strPacket)
				}

				//
				if strings.HasPrefix(strPacket, "GA001") {
					outMoveCharater(id, strPacket)
				}
			}

			//if (strings.HasPrefix(packet, ""))
		},
	}

	err := p.ListenAndServe()
	if err != nil {
		fmt.Println(err)
	}

}

func inputKeyboard() {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("Enter your command: ")
		command, _ := reader.ReadString('\n')
		fmt.Print("Command " + command)

		if command == "GO\n" {
			char := getChararacter("Doelia")
			fmt.Println("write " + char.Name)

			packetConfirm := bytes.NewBufferString("GA001bdVadW")
			packetConfirm.WriteByte(0)
			packetConfirm.WriteString("\n")
			n, _ := char.ConnServer.Write(packetConfirm.Bytes())
			fmt.Println(n)
		}
	}
}

func main() {
	go login()
	go game()
	inputKeyboard()
}