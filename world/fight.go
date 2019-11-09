package world


type Fighter struct {
	Id string
	Name string
	IsMonster bool
	Life int
	Level int
	TeamId int

	IsMe bool
	IsMyTeam bool

	CellId int
	PA int
	PM int
	MaxLife int
}

type Fight struct {
	IdPlayerStarter string
	Fighters []Fighter
}

func (fight Fight) MapId() int {
	return GetPlayer(fight.IdPlayerStarter).MapId
}

func (fight Fight) GetFighter(fighterId string) *Fighter {
	for i, c := range fight.Fighters {
		if c.Id == fighterId || c.Name == fighterId {
			return &fight.Fighters[i]
		}
	}
	return nil
}

func (fighter Fighter) GetPorteeOfBestCast() int {
	return 11
}

func (fight Fight) GetTeamOfFighter(fighterId string) int {
	return fight.GetFighter(fighterId).TeamId
}

func (fight Fight) AreInSameTeam(id string, id2 string) bool {
	return fight.GetTeamOfFighter(id) == fight.GetTeamOfFighter(id2)
}

func UpdateFighter(fight *Fight, fighter Fighter) {
	f := fight.GetFighter(fighter.Id)
	if f != nil {
		f.CellId = fighter.CellId
	} else {
		fight.Fighters = append(fight.Fighters, fighter)
	}
}