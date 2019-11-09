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

	if player.Fight == nil {
		fmt.Println("OnStartTurn: no fight found")
		return
	}

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
			go AutoPlaytTurn(*player)
		}
	}
}

// GJK2|0|1|0|30000|4
func OnJoinFight(player *world.Player, packet string) {
	fmt.Println("OnJoinFight", player.Name)

	player.Fight = &world.Fight{}
	player.Fight.IdPlayerStarter = player.IdCharDofus

	themap := database.GetMap(player.MapId)
	web.SendMap(themap)
}

func OnEndFight(player *world.Player, packet string) {
	fmt.Println("OnEndFight")
	player.Fight = nil
	web.SendCharacters(world.Players)
}


// GTM|-1;0;2;4;2;370;;15|90094963;0;53;6;3;371;;55
func OnFighterUpdateInfos(player *world.Player, packet string) {

	if player.Fight == nil {
		return
	}

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

	if player.Fight == nil {
		return
	}

	idFighterDead := args[3]
	fighterDead := player.Fight.GetFighter(idFighterDead)
	if fighterDead != nil {
		fighterDead.Life = 0
	}

	fmt.Println("fighter is dead", fighterDead)
}
