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

func AddMoveTo(idPlayer string, startMap int, goalMap int) {

	// Remove old if exists
	RemoveMoveTo(idPlayer)

	fmt.Println("AddMoveTo", idPlayer, startMap, goalMap)

	path := world.AStarInWorld(startMap, goalMap)
	fmt.Println("path found", path)

	moveTo := MoveTo{
		IdPlayer: idPlayer,
		pathStack: path,
	}

	MoveToArray = append(MoveToArray, moveTo)
}

func processMoveTo(player world.Player) {

	if player.MapId == 0 || player.CellId == 0 {
		fmt.Println("[moveto] ERROR : position of player invalid", player.Name, player.MapId, player.CellId)
	}

	moveTo := GetMoveTo(player.IdCharDofus)
	if moveTo != nil {

		fmt.Println("[moveto] Player", player.Name, "arrive on map", player.MapId, "cellid=", player.CellId, "path=", moveTo)

		nextMap := 0
		for index, mapID := range moveTo.pathStack {
			if mapID == player.MapId {
				if index == len(moveTo.pathStack) - 1 {
					fmt.Println("[MoveTo]", player.Name, " is arrived on goal map", player.MapId)
					RemoveMoveTo(player.IdCharDofus)
					return
				}
				nextMap = moveTo.pathStack[index + 1]
			}
		}

		if nextMap == 0 {
			fmt.Println("[MoveTo] no next map found from map", player.MapId)
			return
		}

		fmt.Println("[MoveTo]", player.Name, "have to go to", nextMap)

		cell := world.GetCellToGoToMap(player.MapId, nextMap)

		if cell != 0 {

			themap := database.GetMap(player.MapId)
			cellsPath := world.AStar(themap, nil, player.CellId, cell)
			fmt.Println("[MoveTo] path to cell", cell, " is", cellsPath)

			if len(cellsPath) > 0 {
				pathEncoded := world.EncodePath(themap, cellsPath)
				socket.SendMovePacket(*player.Connexion, pathEncoded)
			}
		} else {
			fmt.Println("[MoveTo] cell to go from map", player.MapId, "to map", nextMap, "not found")
		}

	}

}
