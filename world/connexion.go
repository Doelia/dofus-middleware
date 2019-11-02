package world

import (
	"fmt"
	"net"
)

var Connexions []Connexion

type Connexion struct {
	Id string
	Player *Player
	ConnClient net.Conn
	ConnServer net.Conn
}

func AddConnexion(conn Connexion) {
	fmt.Println("Create connexion", conn.Id)
	Connexions = append(Connexions, conn)
}

func GetConnexion(id string) *Connexion {
	for i, c := range Connexions {
		if c.Id == id {
			return &Connexions[i]
		}
	}
	return nil
}

func RemoveConnexion(id string) {
	fmt.Println("Remove connexion " + id)
	var newConnexions []Connexion
	for _, c := range Connexions {
		if c.Id != id {
			newConnexions = append(newConnexions, c)
		}
	}
	Connexions = newConnexions
}
