package world

type Trigger struct {
	FromCellID int
	ToMapId int
	ToCellId int
}

type MapWithTrigger struct {
	MapId int
	Triggers []Trigger
	X int
	Y int
}

var MapWithTriggers []MapWithTrigger

func GetMapWithTriggerWithID(idmap int) MapWithTrigger {
	for _, mapWithTrigger := range MapWithTriggers {
		if mapWithTrigger.MapId == idmap {
			return mapWithTrigger
		}
	}
	return MapWithTrigger{}
}

func GetCellToGoToMap(from_map int, to_map int) int {
	for _, trigger := range GetMapWithTriggerWithID(from_map).Triggers {
		if trigger.ToMapId == to_map {
			return trigger.FromCellID
		}
	}
	return 0
}
