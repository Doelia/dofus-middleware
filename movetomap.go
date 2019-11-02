package main

import (
	"dofusmiddleware/database"
	"dofusmiddleware/socket"
	"dofusmiddleware/world"
	"fmt"
)

type MoveTo struct {
	IdPlayer string
	pathStack []int
}

var MoveToArray []MoveTo

func GetMoveTo(idPlayer string) *MoveTo {
	for _, v := range MoveToArray {
		if v.IdPlayer == idPlayer {
			return &v
		}
	}
	return nil
}

func RemoveMoveTo(idPlayer string) {
	var newArray []MoveTo
	for _, v := range MoveToArray {
		if v.IdPlayer != idPlayer {
			newArray = append(newArray, v)
		}
	}
	MoveToArray = newArray
}

func AddMoveTo(idPlayer string, startMap, goalMap int) {

	moveTo := MoveTo{
		IdPlayer: idPlayer,
		pathStack: world.AStarInWorld(startMap, goalMap)[1:],
	}

	fmt.Println("AddMoveTo", moveTo)

	MoveToArray = append(MoveToArray, moveTo)
}

func OnArriveOnMap(player world.Player) {

	moveTo := GetMoveTo(player.IdCharDofus)
	if moveTo != nil {

		nextMap := moveTo.pathStack[0]
		moveTo.pathStack = moveTo.pathStack[:1]

		if nextMap == player.MapId {
			fmt.Println("[MoveTo]", player.Name, " is arrived on goal map", player.MapId)
			RemoveMoveTo(player.IdCharDofus)
			return
		} else {
			fmt.Println("[MoveTo]", player.Name, "have to", moveTo.pathStack)
			// unstack

			cell := world.GetCellToGoToMap(player.MapId, nextMap)

			themap := database.GetMap(player.MapId)
			cellsPath := world.AStar(themap, player.CellId, cell)
			pathEncoded := world.EncodePath(themap, cellsPath)
			socket.SendMovePacket(*player.Connexion, pathEncoded)
		}
	}

}
