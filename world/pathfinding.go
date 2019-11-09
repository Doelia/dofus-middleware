package world

import (
	"dofusmiddleware/tools"
	"fmt"
)

// cost is the heuristic function. h(n) estimates the cost to reach goal from node n.
func cost(themap Map, cell int, goal int) int {
	return DistanceBetween(themap, cell, goal)
}

func EncodePath(themap Map, cells []int) string {
	encoded := ""

	if len(cells) < 2 {
		return ""
	}

	for i, cell := range cells[1:] { // dont encode the start cell
		orientation := GetDirection(themap, cells[i], cell)
		cellEncoded := encodeOrientedCell(orientation, cell)
		encoded = encoded + cellEncoded
	}

	return encoded
}

func encodeOrientedCell(direction int, cellid int) string {
	return string(tools.EncodeChar(direction)) + string(tools.EncodeChar(cellid / 64)) + string(tools.EncodeChar(cellid % 64))
}

func getNeighbors(themap Map, cell int, withDiagonals bool) []int {
	var cells []int
	c1 := GetCellOfMap(themap, cell + themap.Width)
	c2 := GetCellOfMap(themap, cell - themap.Width)
	c3 := GetCellOfMap(themap, cell - (themap.Width - 1))
	c4 := GetCellOfMap(themap, cell + (themap.Width - 1))

	if c1.Movement > 0 {
		cells = append(cells, c1.CellId)
	}
	if c2.Movement > 0 {
		cells = append(cells, c2.CellId)
	}
	if c3.Movement > 0 {
		cells = append(cells, c3.CellId)
	}
	if c4.Movement > 0 {
		cells = append(cells, c4.CellId)
	}

	if withDiagonals {
		c5 := GetCellOfMap(themap, cell + 1)
		c6 := GetCellOfMap(themap, cell - 1)
		c7 := GetCellOfMap(themap, cell - (themap.Width * 2) + 1)
		c8 := GetCellOfMap(themap, cell + (themap.Width * 2) - 1)

		if c5.Movement > 0 {
			cells = append(cells, c5.CellId)
		}
		if c6.Movement > 0 {
			cells = append(cells, c6.CellId)
		}
		if c7.Movement > 0 {
			cells = append(cells, c7.CellId)
		}
		if c8.Movement > 0 {
			cells = append(cells, c8.CellId)
		}
	}

	return cells
}

func reconstruct_path(cameFrom map[int]int, current int) []int {
	total_path := []int{current}
	current, exists := cameFrom[current]
	for exists {
		total_path = append([]int{current}, total_path...)
		current, exists = cameFrom[current]
	}
	return total_path
}

var Visited []int

// A* finds a path from start to goal.
func AStar(themap Map, start int, goal int, withDiagonals bool) []int {

	Visited = []int{}

	fmt.Println("AStar into map", themap.MapId, "from cell", start, "to cell", goal)

	// The set of discovered nodes that may need to be (re-)expanded.
	// Initially, only the start node is known.
	openSet := []int{start}

	// For node n, cameFrom[n] is the node immediately preceding it on the cheapest path from start to n currently known.
	cameFrom := make(map[int]int)

	// For node n, gScore[n] is the cost of the cheapest path from start to n currently known.
	gScore := make(map[int]int)
	gScore[start] = 0

	// For node n, fScore[n] := gScore[n] + h(n).
	fScore := make(map[int]int)
	fScore[start] = cost(themap, start, goal)

	for len(openSet) != 0 {
		// current: the node in openSet having the lowest fScore[] value
		var current int
		for _, cell := range openSet {
			valCell := fScore[cell]
			valCurrent, okCurrent := fScore[current]

			if !okCurrent || valCell <= valCurrent {
				current = cell
			}
		}

		if current == goal {
			return reconstruct_path(cameFrom, current)
		}

		Visited = append(Visited, current)

		openSet = tools.RemoveIntFromSlice(openSet, current)

		for _, neighbor := range getNeighbors(themap, current, withDiagonals) {
			// d(current,neighbor) is the weight of the edge from current to neighbor
			// tentative_gScore is the distance from start to the neighbor through current
			tentative_gScore := gScore[current] + 1

			if val, ok := gScore[neighbor]; !ok || tentative_gScore < val {
				cameFrom[neighbor] = current
				gScore[neighbor] = tentative_gScore
				fScore[neighbor] = gScore[neighbor] + cost(themap, neighbor, goal)

				exists := false
				for _, c := range openSet {
					if c == neighbor {
						exists = true
					}
				}
				if !exists {
					openSet = append(openSet, neighbor)
				}
			}
		}

	}

	fmt.Println("path not found :(")
	return openSet
}


func GetLastCellFromPath(path string) int {
	lastCell := path[len(path)-2:]
	c1 := lastCell[0]
	c2 := lastCell[1]
	return int(tools.DecodeChar(c1))*64 + int(tools.DecodeChar(c2))
}
