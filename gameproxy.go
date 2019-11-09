package main

import (
	"dofusmiddleware/options"
	"dofusmiddleware/socket"
	"dofusmiddleware/tools"
	"dofusmiddleware/world"
	"fmt"
	"strings"
)


func StartGameProxy() {

	fmt.Println("Start Game proxy on 127.0.0.0:5555")

	p := socket.Server{
		Addr:   "127.0.0.1:5555",
		Target: "52.19.56.159:443",
		ModifyResponse: func(b *[]byte, id string) {

			packets := tools.ExtractPackets(b)
			for _, p := range packets {

				strPacket := string(p)
				strPacket = strPacket[:len(strPacket) - 1] // RemoveIntFromSlice trailing '0' byte

				connexion := world.GetConnexion(id)
				player := connexion.Player

				if options.Options.ShowOutputPackets {
					if player != nil {
						fmt.Println("[" + connexion.Player.Name + "] server->client: " + strPacket)
					} else {
						fmt.Println("[?] server->client: " + strPacket)
					}
				}

				if strings.HasPrefix(string(p), "ALK") {
					OnCharacterEnterInGame(connexion, strPacket)
				}

				if player == nil {
					continue
				}

				if strings.HasPrefix(string(p), "GTS") {
					OnStartTurn(player, strPacket)
				}

				if strings.HasPrefix(string(p), "PIK") {
					OnPopupGroupInvitation(player, strPacket)
				}

				if strings.HasPrefix(string(p), "ERK") {
					go OnPopupExchange(player, strPacket)
				}

				if strings.HasPrefix(string(p), "Gt") {
					go OnFightPopOnMap(player, strPacket)
				}

				if strings.HasPrefix(string(p), "GM") {
					OnSpriteInformation(player, strPacket)
				}
				
				if strings.HasPrefix(string(p), "GJK") {
					OnJoinFight(player, strPacket)
				}
				
				if strings.HasPrefix(string(p), "GE") {
					OnEndFight(player, strPacket)
				}
				
				if strings.HasPrefix(string(p), "GDM") {
					OnMapInfo(player, strPacket)
				}

				if strings.HasPrefix(string(p), "GTM") {
					OnFighterUpdateInfos(player, strPacket)
				}

				if strings.HasPrefix(string(p), "GA0") {
					OnCharacterMove(player, strPacket)
				}

				if strings.HasPrefix(string(p), "GA;103") {
					OnFighterDead(player, strPacket)
				}
			}

		},
		ModifyRequest: func(b *[]byte, id string) {

			bytess := make([]byte, len(*b))
			copy(bytess, *b)
			bytess = bytess[:len(bytess) - 1] // RemoveIntFromSlice trailing '\n' byte
			bytess[len(bytess) - 1] = 0

			packets := tools.ExtractPackets(&bytess)
			for _, p := range packets {

				strPacket := string(p)
				strPacket = strPacket[:len(strPacket) - 1] // RemoveIntFromSlice trailing '0' byte

				connexion := world.GetConnexion(id)
				player := connexion.Player
				if player != nil && options.Options.ShowInputPackets {
					fmt.Println("[" + connexion.Player.Name + "] client->server: " + strPacket)
				}

				if strings.HasPrefix(strPacket, "GA001") {
					OnMoveCharater(player, strPacket)
				}
			}

		},
	}

	err := p.ListenAndServe()
	if err != nil {
		fmt.Println(err)
	}

}
