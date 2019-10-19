package main

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

var characters []Character

type Character struct {
	Name string
	IdCharDofus string
	Id string
}

func getChararacter(id string) *Character {
	for _, c := range characters {
		if c.Id == id {
			return &c
		}
	}
	return nil
}

func login() {

	fmt.Print("Hello world\n");

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
	characters = append(characters, Character{
		Id: id,
		Name: name,
		IdCharDofus: params[0],
	})
}

func startTurn(id string, packet string) {
	splited := strings.Split(packet[3:], "|")
	idCharTurn := splited[0]
	char := getChararacter(id)
	if char.IdCharDofus == idCharTurn {
		fmt.Println("Start turn of " + char.Name)

		cmd := "/Users/stephane/Documents/dev/perso/dofus/" + char.Name + ".sh"
		out := exec.Command("/bin/bash", cmd)
		out.Run()
	}
}

func game() {

	fmt.Print("Hello world\n");

	p := Server{
		Addr:   "127.0.0.1:5555",
		Target: "52.19.56.159:443",
		ModifyResponse: func(b *[]byte, id string) {
			//fmt.Println(*b)

			packets := extractPackets(b)
			for _, p := range packets {
				fmt.Println("[game] server->client: " + string(p))
				if strings.HasPrefix(string(p), "ALK") {
					processPerso(id, string(p))
				}

				if strings.HasPrefix(string(p), "GTS") {
					startTurn(id, string(p))
				}
			}

		},
		ModifyRequest: func(b *[]byte, id string) {
			packet := string(*b)

			fmt.Println("[game] client->server: " + packet)

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