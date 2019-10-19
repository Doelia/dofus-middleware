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
		ModifyResponse: func(b *[]byte) {
			packet := string(*b)
			fmt.Println("[login] server->client: " + packet)
			fmt.Println(*b)
			if strings.Contains(packet, "AXK3413389?ag78d352") {
				fmt.Println("bibly packet transform!")
				by := bytes.NewBufferString("AYK127.0.0.1:5555;8d352")
				//by := bytes.NewBufferString("AXK3413389?ag78d352")
				by.Write([]byte{0})
				*b = by.Bytes()

				fmt.Println(*b)

			}
		},
	}

	err := p.ListenAndServe()
	if err != nil {
		fmt.Println(err)
	}

}

func game() {

	fmt.Print("Hello world\n");

	p := Server{
		Addr:   "127.0.0.1:5555",
		Target: "52.19.56.159:443",
		ModifyResponse: func(b *[]byte) {
			packet := string(*b)

			fmt.Println("[game] server->client: " + packet)
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