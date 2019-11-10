package main


import (
	"dofusmiddleware/database"
	"dofusmiddleware/socket"
	"dofusmiddleware/world"
	"fmt"
	"time"
)

func BotRoutine(p *world.Player) {
	for {
		time.Sleep(time.Duration(5) * time.Second)

		if p == nil {
			break
		}

		if p.OptionAutoStartFight {
			fmt.Println("Execution BotRoutine of", p)
			go SearchNextFight(p)
		}
	}
}

func SearchNextFight(p *world.Player) {

	if !p.OptionAutoStartFight {
		return
	}

	if p.Life < p.MaxLife {
		fmt.Println("[SearchNextFight] Life not full. Wait for it.", p)
		if !p.IsSit {
			go socket.SendSit(*p.Connexion)
		}
		//timeToWait := p.MaxLife - p.Life
		//time.Sleep(time.Duration(timeToWait) * time.Second * 2)
		return
	}

	if p.Fight == nil {
		target, err := p.GetAFigthableEntity()
		if err == nil {
			fmt.Println("[SearchNextFight] target", target)
			themap := database.GetMap(p.MapId)
			path := world.AStar(themap, p.Fight, p.CellId, target.CellId)
			pathEncoded := world.EncodePath(themap, path)
			if pathEncoded != "" {
				socket.SendMovePacket(*p.Connexion, pathEncoded)
			}
		} else {
			fmt.Println("[SearchNexFight] No entity found.")
		}
	} else {
		fmt.Println("[SearchNexFight] In fight.")
	}
}
