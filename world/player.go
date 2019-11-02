package world

import "fmt"

var Players []Player

type Player struct {
	Name string
	IdCharDofus string
	OptionAutoPassTurn bool
	MapId int
	CellId int
	Fight *Fight
	Connexion *Connexion
}

func RemovePlayer(id string) {
	fmt.Println("Remove player " + id)
	var newPlayers []Player
	for _, c := range Players {
		if c.IdCharDofus != id {
			newPlayers = append(newPlayers, c)
		}
	}
	Players = newPlayers
}

func AddPlayer(player Player) {
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
			return &Players[i]
		}
	}
	return nil
}

func GetAConnectedPlayer() *Player {
	for i, c := range Players {
		if len(c.Name) > 0 {
			return &Players[i]
		}
	}
	return nil
}