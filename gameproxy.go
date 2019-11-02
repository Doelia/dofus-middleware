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

func OnJoinFight(id string, packet string) {
	fmt.Println("OnJoinFight")
	getChararacter(id).Fight = &Fight{}
}

func OnEndFight(id string, packet string) {
	fmt.Println("OnEndFight")
	getChararacter(id).Fight = nil
}

func OnMapInfo(id string, packet string) {
	splited := strings.Split(packet, "|")
	idMap := splited[1]
	getChararacter(id).MapId, _ = strconv.Atoi(idMap)
	fmt.Println("map detected", idMap)
}

// [+295 1 0 90069329 Lotahi 9 91^100 1 46 0,0,0,90069375 ffde34 2f8408 295a26 970,96b,96e,6c0, 408 7 3 0 0 0 0 0 20 20 0  ]
// [+170 1 0 -1 236 -2 1212^100 4 a55ee0 ef9f4f -1 0,0,0,0 16 2 3 1]
func OnSpriteInformation(id string, packet string) {

	fmt.Println("Sprite information" + packet)

	entities := strings.Split(packet[3:], "|")
	character := getChararacter(id)

	if character.Fight != nil {

		for _, f := range entities {
			fmt.Println("entity" + f)

			datas := strings.Split(f, ";")

			if len(datas) < 8 {
				fmt.Println("Bad len sprites")
				return
			}

			cellId, _ := strconv.Atoi(datas[0][1:])
			Id, _ := strconv.Atoi(datas[3])
			level, _ := strconv.Atoi(datas[8])

			fmt.Println(datas)
			fmt.Println(Id)

			var fighter Fighter

			if Id < 0 {
				if len(datas) < 15 {
					fmt.Println("Bad len sprites monster")
					return
				}
				teamId, _ := strconv.Atoi(datas[15])
				fighter = Fighter{
					CellId: cellId,
					Id:     datas[3],
					Name:   datas[4],
					Level:  level,
					TeamId: teamId,
				}
			} else {
				if len(datas) < 24 {
					fmt.Println("Bad len sprites player")
					return
				}
				teamId, _ := strconv.Atoi(datas[24])
				fighter = Fighter{
					CellId: cellId,
					Id:     datas[3],
					Name:   datas[4],
					Level:  level,
					TeamId: teamId,
					IsMe: character.Id == datas[3],
				}
			}

			fmt.Println(fighter)
			updateFighter(character.Fight, fighter)
			SendCharacters(Characters)
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

				if strings.HasPrefix(string(p), "GJK") {
					OnJoinFight(id, strPacket)
				}

				if strings.HasPrefix(string(p), "GE") {
					OnEndFight(id, strPacket)
				}

				if strings.HasPrefix(string(p), "GDM") {
					OnMapInfo(id, strPacket)
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
