package dofusmiddleware

import (
	"dofusmiddleware/options"
	"dofusmiddleware/socket"
	"dofusmiddleware/tools"
	"dofusmiddleware/world"
	"fmt"
	"strings"
)


func Game() {

	fmt.Print("New Game proxy")

	p := socket.Server{
		Addr:   "127.0.0.1:5555",
		Target: "52.19.56.159:443",
		ModifyResponse: func(b *[]byte, id string) {
			//fmt.Println(*b)

			packets := tools.ExtractPackets(b)
			for _, p := range packets {

				strPacket := string(p)
				strPacket = strPacket[:len(strPacket) - 1] // Remove trailing '0' byte

				char := world.GetChararacter(id)
				if char != nil && options.Options.ShowOutputPackets {
					fmt.Println("[" + char.Name + "] server->client: " + strPacket)
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

				if strings.HasPrefix(string(p), "Gt") {
					OnFightOpened(id, strPacket)
				}

				//if strings.HasPrefix(string(p), "GM") {
				//	OnSpriteInformation(id, strPacket)
				//}
				//
				//if strings.HasPrefix(string(p), "GJK") {
				//	OnJoinFight(id, strPacket)
				//}
				//
				//if strings.HasPrefix(string(p), "GE") {
				//	OnEndFight(id, strPacket)
				//}
				//
				//if strings.HasPrefix(string(p), "GDM") {
				//	OnMapInfo(id, strPacket)
				//}
				//
				//if strings.HasPrefix(string(p), "GA0") {
				//	OnCharacterMove(id, strPacket)
				//}
			}

		},
		ModifyRequest: func(b *[]byte, id string) {

			bytess := make([]byte, len(*b))
			copy(bytess, *b)
			bytess = bytess[:len(bytess) - 1] // Remove trailing '\n' byte
			bytess[len(bytess) - 1] = 0

			packets := tools.ExtractPackets(&bytess)
			for _, p := range packets {

				strPacket := string(p)
				strPacket = strPacket[:len(strPacket) - 1] // Remove trailing '0' byte

				char := world.GetChararacter(id)
				if char != nil && options.Options.ShowInputPackets {
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
