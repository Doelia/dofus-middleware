package main

import (
	"dofusmiddleware/database"
	"dofusmiddleware/socket"
	"dofusmiddleware/world"
	"fmt"
	"time"
)

func AutoPlayTurn(player world.Player) {

	fight := player.Fight

	if fight != nil {

		me := fight.GetFighter(player.IdCharDofus)
		spell := me.GetBestSpell()
		targetCell := getACellIdCastable(player, spell)

		if targetCell != 0 {
			time.Sleep(time.Duration(300) * time.Millisecond)
			socket.SendCastSpellOnCell(*player.Connexion, spell.IdSpell, targetCell)
		} else {
			time.Sleep(time.Duration(300) * time.Millisecond)
			approachToEnemy(player)
		}

		time.Sleep(time.Duration(300) * time.Millisecond)
		socket.SendPassTurn(*player.Connexion)
	}
}

func getACellIdCastable(player world.Player, spell world.Spell) int {

	fight := player.Fight
	if fight == nil {
		return 0
	}

	me := fight.GetFighter(player.IdCharDofus)

	fmt.Println("[AutoAttack]", me.Name)

	if me.PA < spell.Pa {
		fmt.Println("[AutoAttack] No enough PA", me, spell)
		return 0
	}

	idTarget := getATargetCastable(*fight, player.IdCharDofus, spell)
	fighterTarget := fight.GetFighter(idTarget)
	fmt.Println("[AutoAttack] target chosen", fighterTarget)

	if fighterTarget == nil {
		return 0
	}

	return fighterTarget.CellId

}

func getATargetCastable(fight world.Fight, fighterInt string, spell world.Spell) string {

	themap := database.GetMap(fight.MapId())
	target := ""
	bestDistance := 999

	me := fight.GetFighter(fighterInt)
	for _, fighter := range fight.Fighters {
		if !fight.AreInSameTeam(me.Id, fighter.Id) && fighter.Life > 0 {
			distance := world.DistanceBetween(themap, fighter.CellId, me.CellId)
			if distance <= spell.Portee && distance < bestDistance {
				target = fighter.Id
				bestDistance = distance
			}
		}
	}

	return target
}

func approachToEnemy(player world.Player) string {

	fight := player.Fight
	if fight == nil {
		return ""
	}

	me := fight.GetFighter(player.IdCharDofus)
	idTarget := getATargetToApproach(*fight, player.IdCharDofus)
	fighterTarget := fight.GetFighter(idTarget)

	fmt.Println("[AutoMove] target is", fighterTarget)

	if fighterTarget == nil {
		return ""
	}

	themap := database.GetMap(player.MapId)
	//fmt.Println("[AutoMove] map is", themap.MapId)

	path := world.AStar(themap, me.CellId, fighterTarget.CellId, false)
	//fmt.Println("path from", me.CellId, "to", fighterTarget.CellId, "is", path)

	if len(path) == 0 {
		fmt.Println("[AutoMove] No path found")
		return ""
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

	return pathEncoded

}

func getATargetToApproach(fight world.Fight, fighterInt string) string {

	me := fight.GetFighter(fighterInt)
	for _, fighter := range fight.Fighters {
		if !fight.AreInSameTeam(me.Id, fighter.Id) && fighter.Life > 0 {
			return fighter.Id
		}
	}

	return ""
}
