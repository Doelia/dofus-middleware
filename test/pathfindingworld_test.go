package test

import (
	"dofusmiddleware/database"
	"dofusmiddleware/tools"
	"dofusmiddleware/world"
	"testing"
)

func TestPathFindingWorld(t *testing.T) {
	world.MapWithTriggers = database.GetMapriggers()
	wanted := []int{957,951,952}
	got := world.AStarInWorld(690, 952)
	got = got[1:]
	if !tools.Equal(wanted, got) {
		t.Error("Failed, wanted=", wanted, "path=", got)
	}
}

func TestGetMapFromPosition(t *testing.T) {
	got := database.GetMapIdFromPosition("5,-30")
	wanted := 1868
	if wanted != got {
		t.Error("Failed, wanted=", wanted, "path=", got)
	}
}