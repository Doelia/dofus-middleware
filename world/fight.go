package world


type Fighter struct {
	Id string
	Name string
	IsMonster bool
	Life int
	Level int
	TeamId int

	CellId int
	IsMe bool
	PA int
	PM int
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