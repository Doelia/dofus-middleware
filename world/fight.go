package world


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

func GetFighter(fight *Fight, fighterId string) *Fighter {
	for i, c := range fight.Fighters {
		if c.Id == fighterId || c.Name == fighterId {
			return &fight.Fighters[i]
		}
	}
	return nil
}


func UpdateFighter(fight *Fight, fighter Fighter) {
	f := GetFighter(fight, fighter.Id)
	if f != nil {
		f.CellId = fighter.CellId
	} else {
		fight.Fighters = append(fight.Fighters, fighter)
	}
}