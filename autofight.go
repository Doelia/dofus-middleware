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

func AutoPlaytTurn(player world.Player) {

	time.Sleep(time.Duration(300) * time.Millisecond)
	AutoAttack(player)

	time.Sleep(time.Duration(1000) * time.Millisecond)
	AutoMove(player)

	time.Sleep(time.Duration(200) * time.Millisecond)
	socket.SendPassTurn(*player.Connexion)
}

func AutoAttack(player world.Player) {

	fight := player.Fight
	if fight == nil {
		return
	}

	me := fight.GetFighter(player.IdCharDofus)

	fmt.Println("[AutoAttack]", me.Name)

	if me.PA < me.GetBestSpell().Pa {
		fmt.Println("[AutoAttack] No enough PA", me, me.GetBestSpell())
		return
	}

	idTarget := GetATargetCastable(*fight, player.IdCharDofus)
	fighterTarger := fight.GetFighter(idTarget)
	fmt.Println("[AutoAttack] target chosen", fighterTarger)

	if fighterTarger == nil {
		return
	}

	socket.SendCastSpellOnCell(*player.Connexion, me.GetBestSpell().IdSpell, fighterTarger.CellId)
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