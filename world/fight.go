package world


type Fighter struct {
	Id string
	Name string
	IsMonster bool
	Level int
	TeamId int

	IsMe bool
	IsMyTeam bool

	CellId int
	Life int
	PA int
	PM int
	MaxLife int
}

type Fight struct {
	IdPlayerStarter string
	Fighters []Fighter
}

type Spell struct {
	DebugName string
	IdSpell int
	Portee int
	Pa int
}

const (
	SPELL_LANCER_DE_PIECE = 51
	SPELL_RONCE = 183
	SPELL_FLECHE_MAGIQUE = 161
)

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

func (fighter Fighter) GetBestSpell() Spell {
	return Spell{
		DebugName: "Ronce",
		IdSpell: SPELL_RONCE,
		Portee:  8,
		Pa:      4,
	}
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