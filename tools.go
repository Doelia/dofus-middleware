package main

func remove(slice []int, cell int) []int {
	indexToRemove := 0
	for i, c := range slice {
		if c == cell {
			indexToRemove = i
		}
	}
	return append(slice[:indexToRemove], slice[indexToRemove+1:]...)
}
