package main

import "fmt"

// cost is the heuristic function. h(n) estimates the cost to reach goal from node n.
func cost(themap Map, cell int, goal int) int {
	return distanceBetween(themap, cell, goal)
}

func encodePath(themap Map, cells []int) string {
	encoded := ""

	for i, cell := range cells[1:] {
		orientation := GetDirection(themap, cells[i], cell)
		cellEncoded := encodeOrientedCell(orientation, cell)
		encoded = encoded + cellEncoded
	}

	return encoded
}

func encodeOrientedCell(direction int, cellid int) string {
	return string(encodeChar(direction)) + string(encodeChar(cellid / 64)) + string(encodeChar(cellid % 64))
}

func getNeighbors(themap Map, cell int) []int {
	var cells []int
	c1 := getCellOfMap(themap, cell + 15)
	c2 := getCellOfMap(themap, cell - 15)
	c3 := getCellOfMap(themap, cell - 14)
	c4 := getCellOfMap(themap, cell + 14)

	if c1.Movement == 4 {
		cells = append(cells, c1.CellId)
	}
	if c2.Movement == 4 {
		cells = append(cells, c2.CellId)
	}
	if c3.Movement == 4 {
		cells = append(cells, c3.CellId)
	}
	if c4.Movement == 4 {
		cells = append(cells, c4.CellId)
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

// A* finds a path from start to goal.
func AStar(themap Map, start int, goal int) []int {

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

		openSet = remove(openSet, current)

		for _, neighbor := range getNeighbors(themap, current) {
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


func getLastCellFromPath(path string) int {
	lastCell := path[len(path)-2:]
	c1 := lastCell[0]
	c2 := lastCell[1]
	return int(decodeChar(c1))*64 + int(decodeChar(c2))
}
