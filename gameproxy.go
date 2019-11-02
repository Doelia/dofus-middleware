package dofusmiddleware

import (
	"dofusmiddleware/options"
	"dofusmiddleware/socket"
	"dofusmiddleware/tools"
	"dofusmiddleware/world"
	"fmt"
	"strings"
)


func StartGameProxy() {

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

				character := world.GetChararacter(id)
				if character != nil && options.Options.ShowOutputPackets {
					fmt.Println("[" + character.Name + "] server->client: " + strPacket)
				}

				if strings.HasPrefix(string(p), "ALK") {
					OnCharacterEnterInGame(character, strPacket)
				}

				if strings.HasPrefix(string(p), "GTS") {
					OnStartTurn(character, strPacket)
				}

				if strings.HasPrefix(string(p), "PIK") {
					OnPopupGroupInvitation(character, strPacket)
				}

				if strings.HasPrefix(string(p), "ERK") {
					OnPopupExchange(character, strPacket)
				}

				if strings.HasPrefix(string(p), "Gt") {
					OnFightOpened(character, strPacket)
				}

				if strings.HasPrefix(string(p), "GM") {
					OnSpriteInformation(character, strPacket)
				}
				
				if strings.HasPrefix(string(p), "GJK") {
					OnJoinFight(character, strPacket)
				}
				
				if strings.HasPrefix(string(p), "GE") {
					OnEndFight(character, strPacket)
				}
				
				if strings.HasPrefix(string(p), "GDM") {
					OnMapInfo(character, strPacket)
				}
				
				if strings.HasPrefix(string(p), "GA0") {
					OnCharacterMove(character, strPacket)
				}
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

				character := world.GetChararacter(id)
				if character != nil && options.Options.ShowInputPackets {
					fmt.Println("[" + character.Name + "] client->WebSocket: " + strPacket)
				}

				if strings.HasPrefix(strPacket, "GA001") {
					OnMoveCharater(character, strPacket)
				}
			}

		},
	}

	err := p.ListenAndServe()
	if err != nil {
		fmt.Println(err)
	}

}
