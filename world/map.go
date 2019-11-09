package world

import (
	"fmt"
	"math"
)

type Map struct {
	MapId int
	Cells []Cell
	Width int
	Height int
}

type Cell struct {
	CellId int
	LineOfSight bool
	Movement int
	Interactive bool
	ObjectId int
	Active bool
}

func (m Map) IsATriggerCell(id int) bool {
	return GetCellOfMap(m, id).Movement == 2
}

func GetCellOfMap(themap Map, cellid int) Cell {
	for _, c := range themap.Cells {
		if c.CellId == cellid {
			return c
		}
	}
	fmt.Println("Cant find cellid ", cellid, " of map")
	return Cell{}
}



func GetCellXCoord(mape Map, cellID int) int {
	width := mape.Width
	return (cellID - (width - 1) * GetCellYCoord(mape, cellID)) / width
}

func GetCellYCoord(mape Map, cellID int) int {
	width := mape.Width
	loc5 := cellID / ((width * 2) - 1)
	loc6 := cellID - loc5 * ((width * 2) - 1)
	loc7 := loc6 % width
	return loc5 - loc7
}

func DistanceBetween(mape Map, id1 int, id2 int) int {
	diffX := math.Abs(float64(GetCellXCoord(mape, id1) - GetCellXCoord(mape, id2)))
	diffY := math.Abs(float64(GetCellYCoord(mape, id1) - GetCellYCoord(mape, id2)))

	return int(diffX + diffY)
}

func GetDirection(mape Map, cell1 int, cell2 int) int {

	ListChange := [...] int{
		1,
		mape.Width,
		mape.Width * 2 - 1,
		mape.Width - 1, -1,
		-mape.Width,
		-mape.Width * 2 + 1,
		-(mape.Width - 1),
	}

	Result := cell2 - cell1

	for i := 7; i > -1; i-- {
		if Result == ListChange[i] {
			return i
		}
	}

	ResultX := GetCellXCoord(mape, cell2) - GetCellXCoord(mape, cell1);
	ResultY := GetCellYCoord(mape, cell2) - GetCellYCoord(mape, cell1);

	if ResultX == 0 {
		if ResultY > 0 {
			return 3
		} else {
			return 7
		}
	} else if ResultX > 0 {
			return 1
	} else {
		return 5
	}
}

func OppositeDirection(direction int) int {
	if direction >= 4 {
		return direction - 4
	} else {
		return direction + 4
	}
}

func InLine(mape Map, cell1 int, cell2 int) bool {
	isX := GetCellXCoord(mape, cell1) == GetCellXCoord(mape, cell2);
	isY := GetCellYCoord(mape, cell1) == GetCellYCoord(mape, cell2);
	return isX || isY
}