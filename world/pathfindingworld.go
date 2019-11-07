package world

import (
	"dofusmiddleware/tools"
	"fmt"
	"math"
)

func DistanceBetweenMaps(m1 MapWithTrigger, m2 MapWithTrigger) int {
	diffX := math.Abs(float64(m1.X - m2.X))
	diffY := math.Abs(float64(m1.Y - m2.Y))

	return int(diffX + diffY)
}

// cost is the heuristic function. h(n) estimates the cost to reach goal from node n.
func costWorld(themap int, goal int) int {
	return DistanceBetweenMaps(
		GetMapWithTriggerWithID(themap),
		GetMapWithTriggerWithID(goal),
	)
}

func getNeighborsOfMap(idmap int) []int {
	themap := GetMapWithTriggerWithID(idmap)
	var maps []int
	for _, trigger := range themap.Triggers {
		maps = append(maps, trigger.ToMapId)
	}
	return maps
}

// A* finds a path from start to goal.
func AStarInWorld(start int, goal int) []int {

	fmt.Println("AStarInWorld", start, goal)

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
	fScore[start] = costWorld(start, goal)

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

		openSet = tools.RemoveIntFromSlice(openSet, current)

		for _, neighbor := range getNeighborsOfMap(current) {
			// d(current,neighbor) is the weight of the edge from current to neighbor
			// tentative_gScore is the distance from start to the neighbor through current
			tentative_gScore := gScore[current] + 1

			if val, ok := gScore[neighbor]; !ok || tentative_gScore < val {
				cameFrom[neighbor] = current
				gScore[neighbor] = tentative_gScore
				fScore[neighbor] = gScore[neighbor] + costWorld(neighbor, goal)

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

	fmt.Println("path world not found :(")
	return openSet
}


