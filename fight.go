package main


type Fighter struct {
	Id string
	Name string
	Level int
	TeamId int
	CellId int
	Life int
	IsMe bool
	IsMonster bool
}

type Fight struct {
	Fighters []Fighter
}

func getFighter(fight *Fight, fighterId string) *Fighter {
	for i, c := range fight.Fighters {
		if c.Id == fighterId || c.Name == fighterId {
			return &fight.Fighters[i]
		}
	}
	return nil
}


func updateFighter(fight *Fight, fighter Fighter) {
	f := getFighter(fight, fighter.Id)
	if f != nil {
		f.CellId = fighter.CellId
	} else {
		fight.Fighters = append(fight.Fighters, fighter)
	}
}