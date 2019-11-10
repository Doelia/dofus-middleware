package world

import (
	"encoding/json"
	"fmt"
	"time"
)

const (
	RACE_FECA = 1
	RACE_OSAMODAS = 2
	RACE_ENUTROF = 3
	RACE_SRAM = 4
	RACE_XELOR = 5
	RACE_ECAFLIP = 6
	RACE_ENIRIPSA = 7
	RACE_IOP = 8
	RACE_CRA = 9
	RACE_SADIDA = 10
	RACE_SACRIEUR = 11
	RACE_PANDAWA = 12
)

const (
	MALE = 0
	FEMALE = 1
)

var Players []*Player

type Player struct {
	Name               string
	IdCharDofus        string
	RaceId int
	OptionAutoPassTurn bool
	MapId              int
	CellId             int
	Fight              *Fight
	Connexion          *Connexion
	Life 			int
	MaxLife	int
	IsSit bool

	EntitiesOnSameMap []EntityOnMap

	OptionAutoFight    bool
	OptionAutoStartFight bool
}

// http://jb-serv.forumsactifs.net/t11-id-des-morph
func RaceFromSprite(spriteId int) (raceId int, sexe int) {
	if spriteId >= 10 && spriteId <= 121 {
		return spriteId/10, sexe % 10
	}
	return 0, 0
}

func (p *Player) RegenerateVitaRoutine() {
	for {
		time.Sleep(time.Duration(2) * time.Second)
		if p == nil {
			return
		}
		if p.Fight == nil && p.Life < p.MaxLife {
			if p.IsSit {
				p.Life += 2
			} else {
				p.Life += 1
			}
		}
	}
}

func (p Player) ToJson() ([]byte, error) {
	type player Player
	x := player(p)
	x.Connexion = nil
	return json.Marshal(x)
}

func RemovePlayer(id string) {
	fmt.Println("Remove player " + id)
	var newPlayers []*Player
	for _, c := range Players {
		if c.IdCharDofus != id {
			newPlayers = append(newPlayers, c)
		}
	}
	Players = newPlayers
}

func AddPlayer(player *Player) {
	Players = append(Players, player)
}

func IsOneOfMyPlayer(name string) bool {
	for _, c := range Players {
		if c.Name == name || c.IdCharDofus == name {
			return true
		}
	}
	return false
}

func GetPlayer(search string) *Player {
	for i, c := range Players {
		if c.IdCharDofus == search || c.Name == search {
			return Players[i]
		}
	}
	fmt.Println("[PlayerCollection] cant find player", search)
	return nil
}

func GetAConnectedPlayer() *Player {
	for i, c := range Players {
		if len(c.Name) > 0 {
			return Players[i]
		}
	}
	return nil
}