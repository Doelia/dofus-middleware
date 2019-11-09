package database

import (
	"database/sql"
	"dofusmiddleware/options"
	"dofusmiddleware/tools"
	"dofusmiddleware/world"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"strconv"
	"strings"
)

var database *sql.DB


var mapCache map[int]world.Map

func GetMapIdFromPosition(position string) int {
	var err error

	database, err = sql.Open("sqlite3", options.ConfigSqlLitePath)

	like := position + "%"
	count := 0
	idMap := 0

	if err != nil {
		fmt.Println("err", err)
	} else {
		defer database.Close()

		rows, err2 := database.Query("SELECT id FROM maps where mapPos like ?", like)

		if err2 != nil {
			fmt.Println("err", err2)
		} else {
			for rows.Next() {
				count++
				err = rows.Scan(&idMap)
				if err != nil {
					fmt.Println("err", err)
				}
			}
		}
	}

	if count > 1 {
		fmt.Println("Multiple map for this position", position)
		return 0
	} else if count == 1 {
		fmt.Println("[database] map id for position", position, "is", idMap)
		return idMap
	} else {
		fmt.Println("no map found for this position", position)
		return 0
	}
}

func GetMap(id int) world.Map {

	if mapCache != nil {
		if val, ok := mapCache[id]; ok {
			return val
		}
	} else {
		mapCache = make(map[int]world.Map)
	}

	fmt.Println("load map", id, "from database")

	var err error
	database, err = sql.Open("sqlite3", options.ConfigSqlLitePath)

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
				if (mapdata == "") {
					fmt.Println("Error on map", id, "mapdata is empty")
				}
				themap.Cells = buildCellsFromMapData(mapdata)

				mapCache[themap.MapId] = themap
				return themap
			}
		}

	}

	fmt.Println("Error: map", id, "not found in database")
	return world.Map{}
}

func GetMapriggers() []world.MapWithTrigger {

	var triggers []world.MapWithTrigger

	var err error
	database, err = sql.Open("sqlite3", options.ConfigSqlLitePath)

	if err != nil {
		fmt.Println("err", err)
	} else {
		defer database.Close()

		rows, err2 := database.Query("select from_map, GROUP_CONCAT(from_cell || '|' || to_map || '|' || to_cell), m.mappos from triggers t join maps m on t.from_map=m.id group by from_map")

		if err2 != nil {
			fmt.Println("err", err2)
		} else {
			for rows.Next() {
				cellsString := ""
				coords := ""
				mapWithTrigger := world.MapWithTrigger{}
				err = rows.Scan(&mapWithTrigger.MapId, &cellsString, &coords)
				if err != nil {
					fmt.Println("err", err)
				}

				coordsSplited := strings.Split(coords, ",")
				mapWithTrigger.X, _ = strconv.Atoi(coordsSplited[0])
				mapWithTrigger.Y, _ = strconv.Atoi(coordsSplited[1])

				triggersStr := strings.Split(cellsString, ",")
				for _, triggerStr := range triggersStr {
					triggerStrArgs := strings.Split(triggerStr, "|")
					trigger := world.Trigger{}
					trigger.FromCellID, _ = strconv.Atoi(triggerStrArgs[0])
					trigger.ToMapId, _ = strconv.Atoi(triggerStrArgs[1])
					trigger.ToCellId, _ = strconv.Atoi(triggerStrArgs[2])
					mapWithTrigger.Triggers = append(mapWithTrigger.Triggers, trigger)
				}

				triggers = append(triggers, mapWithTrigger)
			}
		}

	}

	return triggers
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
