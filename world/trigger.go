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

func GetMapWithTriggerWithID(idmap int, maps []MapWithTrigger) MapWithTrigger {
	for _, mapWithTrigger := range maps {
		if mapWithTrigger.MapId == idmap {
			return mapWithTrigger
		}
	}
	return MapWithTrigger{}
}
