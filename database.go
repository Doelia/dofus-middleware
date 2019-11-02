package main

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
)

var database *sql.DB

func getMap(id int) Map {

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
				themap := Map{}
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
	return Map{}
}