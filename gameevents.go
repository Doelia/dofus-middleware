package main

import (
	"bytes"
	"dofusmiddleware/database"
	"dofusmiddleware/options"
	"dofusmiddleware/socket"
	"dofusmiddleware/tools"
	"dofusmiddleware/world"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// ALK708508177|2|90069284;Doelia;51;71;380744;510022;efb57a;30e,1b0e,1b0f,,;0;609;;;
func OnCharacterEnterInGame(connexion *world.Connexion, packet string) {
	splited := strings.Split(packet, "|")
	pr := splited[2]

	params := strings.Split(pr, ";")
	name := params[1]
	idSprite, _ := strconv.Atoi(params[3])
	idRace, _ := world.RaceFromSprite(idSprite)

	player := &(world.Player{
		Name: name,
		IdCharDofus: params[0],
		Connexion: connexion,
		RaceId: idRace,
		OptionAutoFight: true,
	})

	connexion.Player = player

	fmt.Println("Player enter in game : " + name, connexion, player)

	world.AddPlayer(player)
	go BotRoutine(player)
	go player.RegenerateVitaRoutine()
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
func OnFightPopOnMap(player *world.Player, packet string) {
	fmt.Println("[" + player.Name + "] OnFightPopOnMap: " + packet)
	splited := strings.Split(packet[2:], "|")
	startedBy := splited[0]

	if world.IsOneOfMyPlayer(startedBy) {
		if options.Options.AutoJoinFight {
			go func () {
				time.Sleep(time.Duration(tools.RandomBetween(200, 2000)) * time.Millisecond)
				socket.JoinFightCharacter(*player.Connexion, startedBy)
			} ()
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

// GA;1;-1;acHdcV
func OnEntityMove(player *world.Player, packet string) {
	fmt.Println("OnEntityMove", packet)
	if player.Fight == nil {
		splited := strings.Split(packet, ";")
		cellId := world.GetLastCellFromPath(splited[3])
		idEntity, _ := strconv.Atoi(splited[2])
		player.UpdateCellId(idEntity, cellId)
	}
}


func OnMapInfo(player *world.Player, packet string) {
	splited := strings.Split(packet, "|")
	idMap, _ := strconv.Atoi(splited[1])
	fmt.Println("[OnMapInfo] Map detected", idMap)
	themap := database.GetMap(idMap)

	player.MapId = idMap
	player.ClearEntityOnMap()

	web.SendMap(themap)

	processMoveTo(*player)

	if player.OptionAutoFight {
		go SearchNextFight(player)
	}
}

// GA0;1;90069329;ae3hen
func OnCharacterMove(player *world.Player, packet string) {
	splited := strings.Split(packet, ";")

	fmt.Println("OnCharacterMove", splited)

	if len(splited) != 4 {
		fmt.Println("Bad player move packet length", splited)
		return
	}

	path := splited[3]
	idChar := splited[2]

	cellId := world.GetLastCellFromPath(path)

	if player.Fight != nil {
		fighter := player.Fight.GetFighter(idChar)
		if fighter != nil { // Can be an invocation
			fmt.Println("[OnCharacterMoveOnFightMap]: player", fighter, "move to ", cellId)
			fighter.CellId = cellId
		}
	} else {
		player := world.GetPlayer(idChar)
		if player != nil {
			fmt.Println("[OnCharacterMoveOnExplorationMap] : player", player.Name, "move to ", cellId)
			player.IsSit = false
			player.CellId = cellId
		}
	}

	web.SendCharacters(world.Players)
}

// ILS1000 sit
// ILS2000 debout
func OnPlayerChangePosition(me *world.Player, packet string) {
	me.IsSit = packet == "ILS1000"
}

// As10264,7300,10500|2540|0|0|0~0,0,0,0,0,0|85,85|9810,10000|30|100|6,0,0,0,6|3,0,0,0,3|30,0,0,0|0,0,0,0|0,0,0,0|0,0,0,0|0,0,0,0|0,0,0,0|0,0,0,0|1,0,0,0|0,0,0,0|0,0,0,0|0,0,0,0|0,0,0,0|0,0,0,0|0,0,0,0|0,0,0,0|0,0,0,0|0,0,0,0|0,0,0,0|0,0,0,0|0,0,0,0|0,0,0,0|0,0,0,0|0,0,0,0|0,0,0,0|0,0,0,0|0,0,0,0|0,0,0,0|0,0,0,0|0,0,0,0|0,0,0,0|0,0,0,0|0,0,0,0|0,0,0,0|0,0,0,0|0,0,0,0|0,0,0,0|0,0,0,0|0,0,0,0|0,0,0,0|0,0,0,0|20
func OnPlayerStats(me *world.Player, packet string) {
	parts := strings.Split(packet[2:], "|")
	lifePart := strings.Split(parts[5], ",")

	me.Life, _ = strconv.Atoi(lifePart[0])
	me.MaxLife, _ = strconv.Atoi(lifePart[1])
}

// GM [+295 1 0 90069329 Lotahi 9 91^100 1 46 0,0,0,90069375 ffde34 2f8408 295a26 970,96b,96e,6c0, 408 7 3 0 0 0 0 0 20 20 0  ]
// GM [+170 1 0 -1 236 -2 1212^100 4 a55ee0 ef9f4f -1 0,0,0,0 16 2 3 1]
// GM group on exploration map : [+193 5 0 -1 970,970,979 -3 1558^110,1558^100,1560^90 6,4,2 -1,-1,-1 0,0,0,0 -1,-1,-1 0,0,0,0 -1,-1,-1 0,0,0,0 ]
func OnSpriteInformation(me *world.Player, packet string) {

	entities := strings.Split(packet[3:], "|")

	if me.Fight == nil {
		for _, f := range entities {
			datas := strings.Split(f, ";")

			fmt.Println("[OnSpriteInformation/ExplorationMap] entity ", datas)

			if len(datas) > 9 {

				entity := world.BuildEntity(datas)
				fmt.Println("[[OnSpriteInformation/ExplorationMap] entity builded", entity)

				me.AddEntityOnMap(entity)

				if entity.Id > 0 {
					player := world.GetPlayer(strconv.Itoa(entity.Id))
					if player != nil {
						player.CellId = entity.CellId
						fmt.Println("[OnSpriteInformation/ExplorationMap] Update player", player.Name, "cellid on map :", player.CellId)

						if player.IdCharDofus == me.IdCharDofus {
							processMoveTo(*player)
						} else {
							fmt.Println("IsntMe", player, me)
						}
					}
				} else {
					if me.OptionAutoStartFight {
						go func () {
							fmt.Println("[OnCreateFoundOnExplorationMap] cell=", entity.CellId)
							time.Sleep(time.Duration(400) * time.Millisecond)
							SearchNextFight(me)
						} ()
					}
				}


			} else if len(datas) == 1 {
				if strings.HasPrefix(datas[0], "--") {
					fmt.Println("[OnSpriteInformation/ExplorationMap] Depop entity", datas[0])
					idEntity, _ := strconv.Atoi(datas[0][1:])
					me.RemoveEntityOnMap(idEntity)
				}
			}
		}
	} else {
		fmt.Println("[OnSpriteInformation/FightMap] packet", packet)

		for _, f := range entities {
			fmt.Println("[OnSpriteInformation/FightMap] loop entity", f)

			datas := strings.Split(f, ";")

			if len(datas) <= 8 {
				fmt.Println("[OnSpriteInformation/FightMap] Bad len sprites")
				return
			}

			cellId, _ := strconv.Atoi(datas[0][1:])
			Id, _ := strconv.Atoi(datas[3])
			level, _ := strconv.Atoi(datas[8])

			fmt.Println(datas)
			fmt.Println(Id)

			var fighter world.Fighter

			if Id < 0 {
				if len(datas) <= 15 {
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
				idSprit, _ := strconv.Atoi(strings.Split(datas[6], "^")[0])
				IdRace, _ := world.RaceFromSprite(idSprit)

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
					RaceId: IdRace,
					IsMe: me.IdCharDofus == datas[3],
				}
			}

			fmt.Println("[SpriteOnFightMap] final fighter: ", fighter)
			world.UpdateFighter(me.Fight, fighter)
			web.SendCharacters(world.Players)
		}
	}
}
