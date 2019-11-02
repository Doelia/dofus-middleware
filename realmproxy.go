package dofusmiddleware

import (
	"bytes"
	"dofusmiddleware/socket"
	"fmt"
	"strings"
)

func StartRealmProxy() {

	fmt.Print("StartRealmProxy")

	p := socket.Server{
		Addr:   "127.0.0.1:478",
		Target: "34.251.172.139:443",
		ModifyResponse: func(b *[]byte, id string) {
			packet := string(*b)
			fmt.Println("[realm] server->client: " + packet)
			if strings.HasPrefix(packet, "AXK3413389?ag7") {
				token := packet[14:]
				fmt.Println("[realm] bibly packet transform!,  token=" + token)
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
