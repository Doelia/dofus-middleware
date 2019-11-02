package world

import (
	"net"
)

var Characters []Character

type Character struct {
	Name string
	IdCharDofus string
	Id string
	ConnClient net.Conn
	ConnServer net.Conn
	OptionAutoPassTurn bool
	MapId int
	CellId int
	Fight *Fight
}

func IsOneOfMyCharacter(name string) bool {
	for _, c := range Characters {
		if c.Name == name || c.IdCharDofus == name {
			return true
		}
	}
	return false
}

func GetChararacter(search string) *Character {
	for i, c := range Characters {
		if c.Id == search || c.Name == search {
			return &Characters[i]
		}
	}
	return nil
}

func GetAConnectedCharacter() *Character {
	for i, c := range Characters {
		if len(c.Name) > 0 {
			return &Characters[i]
		}
	}
	return nil
}