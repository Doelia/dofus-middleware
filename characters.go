package main

import "net"

var Characters []Character

type Character struct {
	Name string
	IdCharDofus string
	Id string
	ConnClient net.Conn
	ConnServer net.Conn
	OptionAutoPassTurn bool
}

func isOneOfMyCharacter(name string) bool {
	for _, c := range Characters {
		if c.Name == name || c.IdCharDofus == name {
			return true
		}
	}
	return false
}

func getChararacter(search string) *Character {
	for i, c := range Characters {
		if c.Id == search || c.Name == search {
			return &Characters[i]
		}
	}
	return nil
}
