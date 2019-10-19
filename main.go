package main

import (
	"bytes"
	"fmt"
	"strings"
)

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

func processPerso(packet string) {
	splited := strings.Split(packet, "|")
	pr := splited[2]

	params := strings.Split(pr, ";")
	name := params[1]

	fmt.Println("Personnage " + name)
}


func game() {

	fmt.Print("Hello world\n");

	p := Server{
		Addr:   "127.0.0.1:5555",
		Target: "52.19.56.159:443",
		ModifyResponse: func(b *[]byte, id string) {
			packet := string(*b)

			fmt.Println("[game] server->client: " + packet)
			//fmt.Println(*b)

			packets := extractPackets(b)
			for _, p := range packets {
				if strings.HasPrefix(string(p), "ALK") {
					processPerso(string(p))
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