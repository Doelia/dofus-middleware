package main


import (
	"dofusmiddleware/database"
	"dofusmiddleware/options"
	"dofusmiddleware/socket"
	"dofusmiddleware/windowmanagement"
	"dofusmiddleware/world"
	"fmt"
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
