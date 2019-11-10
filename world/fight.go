package world

import "errors"

type Fighter struct {
	Id string
	Name string
	IsMonster bool
	Level int
	TeamId int
	RaceId int

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

// http://oxygene-serveur.xooit.fr/t7-Liste-ID-des-sorts-dofus.htm
const (
	SPELL_LANCER_DE_PIECE = 51
	SPELL_RONCE = 183
	SPELL_FLECHE_MAGIQUE = 161
	SPELL_ATTAQUE_NATURELLE = 3
	SPELL_MOT_INTERDIT = 125
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

func (fight Fight) HasEntityOnCellId(cellid int) bool {
	for _, f := range fight.Fighters {
		if f.CellId == cellid {
			return true
		}
	}
	return false
}

func (fighter Fighter) GetBestAttackSpell() (Spell, error) {
	if fighter.RaceId == RACE_SADIDA {
		return Spell{
			DebugName: "Ronce",
			IdSpell: SPELL_RONCE,
			Portee:  8,
			Pa:      4,
		}, nil
	}
	if fighter.RaceId == RACE_CRA {
		return Spell{
			DebugName: "Flèche magique",
			IdSpell: SPELL_FLECHE_MAGIQUE,
			Portee:  11,
			Pa:      4,
		}, nil
	}
	if fighter.RaceId == RACE_ENUTROF {
		return Spell{
			DebugName: "Lancer de pièces",
			IdSpell: SPELL_LANCER_DE_PIECE,
			Portee:  12,
			Pa:      2,
		}, nil
	}
	if fighter.RaceId == RACE_FECA {
		return Spell{
			DebugName: "Attaque naturelle",
			IdSpell: SPELL_ATTAQUE_NATURELLE,
			Portee:  7,
			Pa:      4,
		}, nil
	}
	if fighter.RaceId == RACE_ENIRIPSA {
		return Spell{
			DebugName: "Attaque naturelle",
			IdSpell: SPELL_MOT_INTERDIT,
			Portee:  4,
			Pa:      4,
		}, nil
	}

	return Spell{}, errors.New("No cast found for class")
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