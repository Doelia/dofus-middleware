package main

import (
	"bytes"
	"dofusmiddleware/database"
	"dofusmiddleware/options"
	"dofusmiddleware/socket"
	"dofusmiddleware/windowmanagement"
	"dofusmiddleware/world"
	"fmt"
	"strconv"
	"strings"
	"time"
)

func OnCharacterEnterInGame(connexion *world.Connexion, packet string) {
	splited := strings.Split(packet, "|")
	pr := splited[2]

	params := strings.Split(pr, ";")
	name := params[1]

	player :=  &(world.Player{
		Name: name,
		IdCharDofus: params[0],
		Connexion: connexion,
	})

	connexion.Player = player

	fmt.Println("Player enter in game : " + name, connexion, player)

	world.AddPlayer(player)
}

func OnStartTurn(player *world.Player, packet string) {
	splited := strings.Split(packet[3:], "|")
	idCharTurn := splited[0]
	if player.IdCharDofus == idCharTurn {
		fmt.Println("Start turn of " + player.Name)
		if options.Options.FocusWindowOnCharacterTurn {
			go windowmanagement.SwitchToCharacter(player.Name)
		}
		if player.OptionAutoPassTurn {
			fmt.Println("Pass turn of " + player.Name)
			time.Sleep(time.Duration(200) * time.Millisecond)
			socket.SendPassTurn(*player.Connexion)
		}
	}
}

// PIKDoelia|Lotahi
func OnPopupGroupInvitation(player *world.Player, packet string) {
	splited := strings.Split(packet[3:], "|")
	inviter := splited[0]
	invited := splited[1]

	fmt.Println(inviter + " " + invited + " " + player.Name)

	// Im invited
	if invited == player.Name {
		if world.IsOneOfMyPlayer(inviter) {
			fmt.Println("Im ("+ invited +") invited to join "+ inviter +" group's")
			socket.SendConfirmAction(*player.Connexion)
		}
	}
}

//  ERK90069329|90069284|1
func OnPopupExchange(player *world.Player, packet string) {
	splited := strings.Split(packet[3:], "|")
	inviter := splited[0]
	invited := splited[1]

	fmt.Println(inviter + " " + invited + " " + player.Name)

	// Im invited
	if invited == player.IdCharDofus {
		if world.IsOneOfMyPlayer(inviter) {
			fmt.Println("Im ("+ invited +") invited to exchange with "+ inviter)
			packetConfirm := bytes.NewBufferString("EA")
			packetConfirm.WriteByte(0)
			packetConfirm.WriteString("\n")
			_, _ = player.Connexion.ConnServer.Write(packetConfirm.Bytes())
		}
	}
}

// Gt90069329|+90069329;Lotahi;44
func OnFightOpened(player *world.Player, packet string) {
	fmt.Println("[" + player.Name + "] OnFightOpened: " + packet)
	splited := strings.Split(packet[2:], "|")
	startedBy := splited[0]

	if world.IsOneOfMyPlayer(startedBy) {
		if options.Options.AutoJoinFight {
			go socket.JoinFightCharacter(*player.Connexion, startedBy)
			if options.Options.AutoReadyFight {
				go socket.ReadyFightCharacter(*player.Connexion)
			}
		}
	}
}

// GA001fc4 GA001[move]
func OnMoveCharater(player *world.Player, packet string) {
	if options.Options.DispatchMoves {
		counter := 0
		for _, c := range world.Players {
			if player.IdCharDofus != player.IdCharDofus {
				counter = counter + 1
				fmt.Println(counter)
				go socket.MoveChar(*c.Connexion, packet, counter)
			}
		}
	}
}

// GJK2|0|1|0|30000|4
func OnJoinFight(player *world.Player, packet string) {
	fmt.Println("OnJoinFight")

	player.Fight = &world.Fight{}

	themap := database.GetMap(player.MapId)
	web.SendMap(themap)
}

func OnEndFight(player *world.Player, packet string) {
	fmt.Println("OnEndFight")
	player.Fight = nil
	web.SendCharacters(world.Players)
}

func OnMapInfo(player *world.Player, packet string) {
	splited := strings.Split(packet, "|")
	idMap, _ := strconv.Atoi(splited[1])
	player.MapId = idMap
	fmt.Println("map detected", idMap)
	fmt.Println("player edited", player)
	fmt.Println("player in collection", world.GetPlayer(player.Name))
	themap := database.GetMap(idMap)
	web.SendMap(themap)
}

// GA0;1;90069329;ae3hen
func OnCharacterMove(player *world.Player, packet string) {
	splited := strings.Split(packet, ";")

	fmt.Println("OnCharacterMove", splited)

	if len(splited) != 4 {
		fmt.Println("Bad player move packet length", splited)
		return;
	}

	path := splited[3]
	idChar := splited[2]

	if player.Fight != nil {
		cellId := world.GetLastCellFromPath(path)
		fighter := world.GetFighter(player.Fight, idChar)
		fmt.Println("Fight: player", fighter, "move to ", cellId)
		fighter.CellId = cellId
		web.SendCharacters(world.Players)
	}
}

// GM [+295 1 0 90069329 Lotahi 9 91^100 1 46 0,0,0,90069375 ffde34 2f8408 295a26 970,96b,96e,6c0, 408 7 3 0 0 0 0 0 20 20 0  ]
// GM [+170 1 0 -1 236 -2 1212^100 4 a55ee0 ef9f4f -1 0,0,0,0 16 2 3 1]
func OnSpriteInformation(player *world.Player, packet string) {

	fmt.Println("Sprite information" + packet)

	entities := strings.Split(packet[3:], "|")

	if player.Fight != nil {

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

			var fighter world.Fighter

			if Id < 0 {
				if len(datas) < 15 {
					fmt.Println("Bad len sprites monster")
					return
				}
				teamId, _ := strconv.Atoi(datas[15])
				fighter = world.Fighter{
					IsMonster: true,
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
				fighter = world.Fighter{
					IsMonster: false,
					CellId: cellId,
					Id:     datas[3],
					Name:   datas[4],
					Level:  level,
					TeamId: teamId,
					IsMe: player.IdCharDofus == datas[3],
				}
			}

			fmt.Println("fighter sprite", fighter)
			world.UpdateFighter(player.Fight, fighter)
			web.SendCharacters(world.Players)
		}
	}
}
