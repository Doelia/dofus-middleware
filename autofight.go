package main

import (
	"dofusmiddleware/database"
	"dofusmiddleware/socket"
	"dofusmiddleware/world"
	"fmt"
)

func GetATargetCastable(fight world.Fight, fighterInt string) string {

	themap := database.GetMap(fight.MapId())

	me := fight.GetFighter(fighterInt)
	for _, fighter := range fight.Fighters {
		if !fight.AreInSameTeam(me.Id, fighter.Id) && fighter.Life > 0 {
			distance := world.DistanceBetween(themap, fighter.CellId, me.CellId)
			if distance <= me.GetPorteeOfBestCast() {
				return fighter.Id
			}
		}
	}

	return ""
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

func AutoAttack(player world.Player) {

	fight := player.Fight
	if fight == nil {
		return
	}

	me := fight.GetFighter(player.IdCharDofus)

	fmt.Println("auto attack", me)

	if me.PA < 4 {
		fmt.Println("No enough PA")
		return
	}

	idTarget := GetATargetCastable(*fight, player.IdCharDofus)
	fighterTarger := fight.GetFighter(idTarget)
	fmt.Println("Target is", fighterTarger)

	if fighterTarger == nil {
		fmt.Println("no target found", fighterTarger)
		return
	}

	fmt.Println("cast spell")
	socket.SendCastSpellOnCell(*player.Connexion, 161, fighterTarger.CellId)
}

func AutoMove(player world.Player) {

	fight := player.Fight
	if fight == nil {
		return
	}

	me := fight.GetFighter(player.IdCharDofus)

	fmt.Println("auto move", me)

	idTarget := GetATargetToJoin(*fight, player.IdCharDofus)

	if idTarget == "" {
		fmt.Println("no target found to join")
		return
	}

	fighterTarger := fight.GetFighter(idTarget)

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
		socket.SendMovePacket(*player.Connexion, pathEncoded)
	} else {
		fmt.Println("no path to walk")
	}
}