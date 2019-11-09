package main

import (
	"dofusmiddleware/database"
	"dofusmiddleware/socket"
	"dofusmiddleware/world"
	"fmt"
	"time"
)

func GetATargetCastable(fight world.Fight, fighterInt string) string {

	themap := database.GetMap(fight.MapId())
	target := ""
	bestDistance := 999

	me := fight.GetFighter(fighterInt)
	for _, fighter := range fight.Fighters {
		if !fight.AreInSameTeam(me.Id, fighter.Id) && fighter.Life > 0 {
			distance := world.DistanceBetween(themap, fighter.CellId, me.CellId)
			if distance <= me.GetBestSpell().Portee && distance < bestDistance {
				target = fighter.Id
				bestDistance = distance
			}
		}
	}

	return target
}

func GetATargetToJoin(fight world.Fight, fighterInt string) string {

	me := fight.GetFighter(fighterInt)
	for _, fighter := range fight.Fighters {
		if !fight.AreInSameTeam(me.Id, fighter.Id) && fighter.Life > 0 {
			return fighter.Id
		}
	}

	return ""
}

func AutoReady(player world.Player) {
	time.Sleep(time.Duration(500) * time.Millisecond)
	socket.ReadyFightCharacter(*player.Connexion)
}

func AutoPlaytTurn(player world.Player) {

	fight := player.Fight

	if fight != nil {

		targetCell := AutoAttack(player)
		me := fight.GetFighter(player.IdCharDofus)

		if targetCell != 0 {
			time.Sleep(time.Duration(300) * time.Millisecond)
			socket.SendCastSpellOnCell(*player.Connexion, me.GetBestSpell().IdSpell, targetCell)
		} else {
			time.Sleep(time.Duration(300) * time.Millisecond)
			AutoMove(player)
		}

		time.Sleep(time.Duration(300) * time.Millisecond)
		socket.SendPassTurn(*player.Connexion)
	}
}

func AutoAttack(player world.Player) int {

	fight := player.Fight
	if fight == nil {
		return 0
	}

	me := fight.GetFighter(player.IdCharDofus)

	fmt.Println("[AutoAttack]", me.Name)

	if me.PA < me.GetBestSpell().Pa {
		fmt.Println("[AutoAttack] No enough PA", me, me.GetBestSpell())
		return 0
	}

	idTarget := GetATargetCastable(*fight, player.IdCharDofus)
	fighterTarger := fight.GetFighter(idTarget)
	fmt.Println("[AutoAttack] target chosen", fighterTarger)

	if fighterTarger == nil {
		return 0
	}

	return fighterTarger.CellId

}

func AutoMove(player world.Player) {

	fight := player.Fight
	if fight == nil {
		return
	}

	me := fight.GetFighter(player.IdCharDofus)

	fmt.Println("auto move", me)

	idTarget := GetATargetToJoin(*fight, player.IdCharDofus)
	fighterTarger := fight.GetFighter(idTarget)

	fmt.Println("[AutoMove] target is", fighterTarger)

	if fighterTarger == nil {
		return
	}

	themap := database.GetMap(player.MapId)
	//fmt.Println("[AutoMove] map is", themap.MapId)

	path := world.AStar(themap, me.CellId, fighterTarger.CellId, false)
	//fmt.Println("path from", me.CellId, "to", fighterTarger.CellId, "is", path)

	if len(path) == 0 {
		fmt.Println("[AutoMove] No path found")
		return
	}

	path = path[:len(path)-1] // Remove last cell (is the monster cell)
	sizeMaxPath := me.PM + 1
	if len(path) > sizeMaxPath {
		path = path[:sizeMaxPath]
	}

	fmt.Println("[AutoMove] path to walk is", path, "(", me.PM, " PM)")

	pathEncoded := world.EncodePath(themap, path)
	if pathEncoded != "" {
		socket.SendMovePacket(*player.Connexion, pathEncoded)
	}
}

func OnCreateFoundOnExplorationMap(player *world.Player, cellId int) {
	fmt.Println("[OnCreateFoundOnExplorationMap] cell=", cellId)

	time.Sleep(time.Duration(400) * time.Millisecond)

}

func SearchNextFight(p *world.Player) {

	fmt.Println("[SearchNextFight] Start...")
	time.Sleep(time.Duration(1000) * time.Millisecond)
	if !p.OptionAutoFight {
		return
	}

	if p.Life < p.MaxLife {
		fmt.Println("[SearchNextFight] Life not full. Wait for it.", p)
		go socket.SendSit(*p.Connexion)
		timeToWait := p.MaxLife - p.Life
		time.Sleep(time.Duration(timeToWait) * time.Second)
	}

	if p.Fight == nil {
		target, err := p.GetAFigthableEntity()
		if err == nil {
			fmt.Println("[SearchNextFight] target", target)
			themap := database.GetMap(p.MapId)
			path := world.AStar(themap, p.CellId, target.CellId, true)
			pathEncoded := world.EncodePath(themap, path)
			if pathEncoded != "" {
				socket.SendMovePacket(*p.Connexion, pathEncoded)
			}
		} else {
			fmt.Println("[SearchNexFight] No entity. Wait for next.")
			SearchNextFight(p)
		}
	} else {
		fmt.Println("[SearchNexFight] In fight.")
	}
}