package world

import (
	"encoding/json"
	"fmt"
	"time"
)

const (
	CLASS_IOP = 3
	CLASS_FECA = 4,
)

var Players []*Player

type Player struct {
	Name               string
	IdCharDofus        string
	Class int
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