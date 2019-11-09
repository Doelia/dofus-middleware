package main

import (
	"dofusmiddleware/world"
)

func GetATarget(fight world.Fight, fighterInt string) string {

	me := fight.GetFighter(fighterInt)
	for _, fighter := range fight.Fighters {
		if !fight.AreInSameTeam(me.Id, fighter.Id) && fighter.Life > 0 {
			return fighter.Id
		}
	}

	return ""
}

