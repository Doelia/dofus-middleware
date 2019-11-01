package main

import (
	"fmt"
)

type Map struct {
	MapId int
	Cells []Cell
}

type Cell struct {
	CellId int
	LineOfSight bool
	Movement int
	Interactive bool
	ObjectId int
	Active bool
}


func buildCell(cellData string) Cell {
	fmt.Println(cellData)
	cell := Cell{}
	var data = make([]byte, 10)
	for i := 0; i < len(cellData); i++  {
		char := cellData[i]
		data[i] = decodeChar(char)
	}

	fmt.Println(data)

	cell.LineOfSight = data[0] & 1 == 1
	cell.Movement = int((data[2] & 56) >> 3)
	cell.ObjectId = int(((data[0] & 2) << 12) + ((data[7] & 1) << 12) + (data[8] << 6) + data[9])
	cell.Active = (data[0] & 32) >> 5 == 1

	return cell
}

func buildCellsFromMapData(mapData string) []Cell {
	var cells []Cell
	for i := 0; i < len(mapData); i+=10  {
		cell := buildCell(mapData[i:i+10])
		cell.CellId = i / 10
		cells = append(cells, cell)
	}
	return cells
}