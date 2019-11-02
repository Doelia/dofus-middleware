package database

import (
	"database/sql"
	"dofusmiddleware/tools"
	"dofusmiddleware/world"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
)

var database *sql.DB

func GetMap(id int) world.Map {

	fmt.Println("load map", id, "from database")

	var err error
	database, err = sql.Open("sqlite3", "/Users/stephane/Desktop/dofus.sqllite")

	if err != nil {
		fmt.Println("err", err)
	} else {
		defer database.Close()

		rows, err2 := database.Query("SELECT id, width, heigth, mapData FROM maps where id=?", id)

		if err2 != nil {
			fmt.Println("err", err2)
		} else {
			for rows.Next() {
				themap := world.Map{}
				mapdata := ""
				err = rows.Scan(&themap.MapId, &themap.Width, &themap.Height, &mapdata)
				if err != nil {
					fmt.Println("err", err)
				}
				themap.Cells = buildCellsFromMapData(mapdata)
				return themap
			}
		}

	}

	fmt.Println("Error: map", id, "not found in database")
	return world.Map{}
}

func buildCellsFromMapData(mapData string) []world.Cell {
	var cells []world.Cell
	for i := 0; i < len(mapData); i+=10  {
		cell := buildCell(mapData[i:i+10])
		cell.CellId = i / 10
		cells = append(cells, cell)
	}
	return cells
}

func buildCell(cellData string) world.Cell {
	cell := world.Cell{}
	var data = make([]byte, 10)
	for i := 0; i < len(cellData); i++  {
		char := cellData[i]
		data[i] = tools.DecodeChar(char)
	}

	cell.LineOfSight = data[0] & 1 == 1
	cell.Movement = int((data[2] & 56) >> 3)
	cell.ObjectId = int(((data[0] & 2) << 12) + ((data[7] & 1) << 12) + (data[8] << 6) + data[9])
	cell.Active = (data[0] & 32) >> 5 == 1

	return cell
}
