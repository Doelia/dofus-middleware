package dofusmiddleware

import (
	"bytes"
	"dofusmiddleware/socket"
	"fmt"
	"strings"
)

func Login() {

	fmt.Print("Hello world\n")

	p := socket.Server{
		Addr:   "127.0.0.1:478",
		Target: "34.251.172.139:443",
		ModifyResponse: func(b *[]byte, id string) {
			packet := string(*b)
			fmt.Println("[login] server->client: " + packet)
			if strings.HasPrefix(packet, "AXK3413389?ag7") {
				token := packet[14:]
				fmt.Println("bibly packet transform!,  token=" + token)
				by := bytes.NewBufferString("AYK127.0.0.1:5555;" + token)
				*b = by.Bytes()
				//fmt.Println(*b)
			}
		},
		ModifyRequest: func(b *[]byte, id string) {
			packet := string(*b)
			fmt.Println("[login] client->server: " + packet)
		},
	}

	err := p.ListenAndServe()
	if err != nil {
		fmt.Println(err)
	}

}
