package main

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
	"time"
)

func OnCharacterEnterInGame(id string, packet string) {
	splited := strings.Split(packet, "|")
	pr := splited[2]

	params := strings.Split(pr, ";")
	name := params[1]

	fmt.Println("Character enter in game : " + name)
	var char = getChararacter(id)
	char.Name = name
	char.IdCharDofus = params[0]
}

func OnStartTurn(id string, packet string) {
	splited := strings.Split(packet[3:], "|")
	idCharTurn := splited[0]
	char := getChararacter(id)
	if char.IdCharDofus == idCharTurn {
		fmt.Println("Start turn of " + char.Name)
		if Options.FocusWindowOnCharacterTurn {
			go SwitchToCharacter(char.Name)
		}
		if char.OptionAutoPassTurn {
			fmt.Println("Pass turn of " + char.Name)
			time.Sleep(time.Duration(200) * time.Millisecond)
			packetConfirm := bytes.NewBufferString("Gt")
			packetConfirm.WriteByte(0)
			packetConfirm.WriteString("\n")
			_, _ = char.ConnServer.Write(packetConfirm.Bytes())
		}
	}
}

// PIKDoelia|Lotahi
func OnPopupGroupInvitation(id string, packet string) {
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
func OnPopupExchange(id string, packet string) {
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

// Gt90069329|+90069329;Lotahi;44
func OnFightOpened(id string, packet string) {
	char := getChararacter(id)
	fmt.Println("[" + char.Name + "] OnFightOpened: " + packet)
	splited := strings.Split(packet[2:], "|")
	startedBy := splited[0]

	if isOneOfMyCharacter(startedBy) {
		if Options.AutoJoinFight {
			go joinFightCharacter(*char, startedBy)
			if Options.AutoReadyFight {
				go readyFightCharacter(*char)
			}
		}
	}
}


func OnMoveCharater(id string, packet string) {

	path := packet
	getChararacter(id).CellId = decodePath(path)
	SendCharacters(Characters)

	if Options.DispatchMoves {
		counter := 0
		for _, c := range Characters {
			if c.Name != "" && c.Id != id {
				counter = counter + 1
				fmt.Println(counter)
				go moveChar(c, packet, counter)
			}
		}
	}
}

// monster : GM|
// +339; 0
// 1; 1 orientation
// 0; 2 bonus value ?
// -1; 3 ID
// 490; 4 name
// -2; 5 type id
// 1212^100; 6 gfxid
// 3; 7 grade number
// 448051; 8 couleurs
// f9f9a5; 9 accesoires
// -1; 10
// 0,0,0,0; 11
// 14; 12
// 2; 13
// 3; 14
// 1 15
//
// +324; 0
// 1; 1
// 0; 2
// -2; 3 ID
// 491; 4 name
// -2; 5 type id
// 1212^100; 6 gfxid
// 3; 7 grade number
// -1; 8
// -1; 9
// -1; 10
// 0,0,0,0; 11
// 14; 12
// 2; 13
// 3; 14
// 1 15

// character : GM|
// +119; 0 cell id
// 1; 1 orientation
// 0; 2 bonus value, not used on player
// 90069329; 3
// Lotahi; 4 name
// 9; 5 race
// 91^100; 6 gfxid
// 1; 7 sexe
// 46; 8 level
// 0,0,0,90069375; 9 alignement ?
// ffde34; 10 couleurs
// 2f8408; 11 accsoires
// 295a26; 12
// 970,96b,96e,6c0,; 13
// 408; 14 vie
// 7; 15 PA
// 3; 16 PM
// 0; 17 RES neutre
// 0; 18 RES terre
// 0; 19 res feu
// 0; 20 res eau
// 0; 21 res air
// 20; 22 res PA
// 20; 23 res PM
// 0; 24 num de team
// ; 25
// 26

func OnSpriteInformation(id string, packet string) {

	fmt.Println("Sprite information" + packet)

	entities := strings.Split(packet[3:], "|")

	if CurrentFight != nil {

		for _, f := range entities {
			fmt.Println("entity" + f)

			datas := strings.Split(f, ";")

			cellId, _ := strconv.Atoi(datas[0][1:])
			Id, _ := strconv.Atoi(datas[3])
			level, _ := strconv.Atoi(datas[8])

			var fighter Fighter

			if Id > 0 {
				teamId, _ := strconv.Atoi(datas[15])
				fighter = Fighter{
					CellId: cellId,
					Id:     datas[3],
					Name:   datas[4],
					Level:  level,
					TeamId: teamId,
				}
			} else {
				teamId, _ := strconv.Atoi(datas[15])
				fighter = Fighter{
					CellId: cellId,
					Id:     datas[3],
					Name:   datas[4],
					Level:  level,
					TeamId: teamId,
				}
			}

			updateFighter(*CurrentFight, fighter)
			SendFight(*CurrentFight)
		}
	}
}

func game() {

	fmt.Print("New Game proxy")

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
				if char != nil && Options.ShowOutputPackets {
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

				if strings.HasPrefix(string(p), "GM") {
					OnSpriteInformation(id, strPacket)
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
				if char != nil && Options.ShowInputPackets {
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
