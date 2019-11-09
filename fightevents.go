package main


import (
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

func OnStartTurn(player *world.Player, packet string) {
	splited := strings.Split(packet[3:], "|")
	idCharTurn := splited[0]
	if player.IdCharDofus == idCharTurn {
		fmt.Println("Start of my turn : " + player.Name)
		if options.Options.FocusWindowOnCharacterTurn {
			go windowmanagement.SwitchToCharacter(player.Name)
		}
		if player.OptionAutoPassTurn {
			fmt.Println("Pass turn of " + player.Name)
			time.Sleep(time.Duration(200) * time.Millisecond)
			socket.SendPassTurn(*player.Connexion)
		}
		if player.OptionAutoFight {
			fight := player.Fight
			me := fight.GetFighter(player.IdCharDofus)

			idTarget := GetATarget(*fight, player.IdCharDofus)
			fighterTarger := fight.GetFighter(idTarget)
			fmt.Println("Target is", fighterTarger)

			if fighterTarger == nil {
				fmt.Println("no target found", fighterTarger)
				return
			}

			fmt.Println("cast spell")
			time.Sleep(time.Duration(1000) * time.Millisecond)
			socket.SendCastSpellOnCell(*player.Connexion, 161, fighterTarger.CellId)

			themap := database.GetMap(player.MapId)
			fmt.Println("map is", themap.MapId)

			path := world.AStar(themap, me.CellId, fighterTarger.CellId, false)
			fmt.Println("path from", me.CellId, "to", fighterTarger.CellId, "is", path)

			if len(path) == 0 {
				fmt.Println("No path found")
				return
			}

			path = path[:len(path)-1] // Remove last cell (is the monster cell)
			fmt.Println("path1", path)
			sizeMaxPath := me.PM + 1
			if len(path) > sizeMaxPath {
				path = path[:sizeMaxPath]
			}

			fmt.Println("path to walk is", path, "(", me.PM, " PM)")

			pathEncoded := world.EncodePath(themap, path)
			if pathEncoded != "" {
				time.Sleep(time.Duration(1000) * time.Millisecond)
				socket.SendMovePacket(*player.Connexion, pathEncoded)
			} else {
				fmt.Println("no path to walk")
			}

			time.Sleep(time.Duration(200) * time.Millisecond)
			socket.SendPassTurn(*player.Connexion)


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


// GTM|90069284;0;299;6;3;252;;305|90069329;0;418;7;3;267;;418
func OnFighterUpdateInfos(player *world.Player, packet string) {
	fightersPackets := strings.Split(packet[3:], "|")
	for _, fighterPacket := range fightersPackets {

		args := strings.Split(fighterPacket, ";")
		if len(args) < 8 {
			fmt.Println("OnFighterUpdateInfos bad fighter len", args)
			continue
		}

		idFighter := args[0]
		life, _ := strconv.Atoi(args[2])
		PA, _ := strconv.Atoi(args[3])
		PM, _ := strconv.Atoi(args[4])
		CellId, _ := strconv.Atoi(args[5])
		maxLife, _ := strconv.Atoi(args[7])

		fighter := player.Fight.GetFighter(idFighter)
		if fighter == nil {
			fmt.Println("OnFighterUpdateInfos fighter not found", idFighter)
			continue
		}
		fighter.Life = life
		fighter.PA = PA
		fighter.PM = PM
		fighter.MaxLife = maxLife
		fighter.CellId = CellId

		fmt.Println("Update fighter stats", fighter)
	}

	web.SendCharacters(world.Players)
	fmt.Println("OnFighterUpdateInfos")
}

// GA;103;90069329;-2
func OnFighterDead(player *world.Player, packet string) {

	args := strings.Split(packet, ";")
	if len(args) != 4 {
		fmt.Println("OnFighterDead : bad length packet", args)
		return
	}

	idFighterDead := args[3]
	fighterDead := player.Fight.GetFighter(idFighterDead)
	fighterDead.Life = 0

	fmt.Println("fighter is dead", fighterDead)
}
