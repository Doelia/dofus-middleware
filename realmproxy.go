package main

import (
	"bytes"
	"dofusmiddleware/socket"
	"fmt"
	"strings"
)

func StartRealmProxy() {

	fmt.Println("Start Realm proxy on 127.0.0.1:9000")

	p := socket.Server{
		Addr:   "127.0.0.1:9000",
		Target: "34.251.172.139:443",
		ModifyResponse: func(b *[]byte, id string) {
			packet := string(*b)
			fmt.Println("[realm] server->client: " + packet)
			if strings.HasPrefix(packet, "AXK3413389?ag7") {
				token := packet[14:]
				fmt.Println("[realm] bibly packet transform!,  token=" + token)
				by := bytes.NewBufferString("AYK127.0.0.1:5555;" + token)
				*b = by.Bytes()
			}
		},
		ModifyRequest: func(b *[]byte, id string) {
			packet := string(*b)
			fmt.Println("[realm] client->server: " + packet)
		},
	}

	err := p.ListenAndServe()
	if err != nil {
		fmt.Println(err)
	}

}
