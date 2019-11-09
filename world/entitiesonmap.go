package world

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type EntityOnMap struct {
	CellId int
	Id int
	levels []int
}

func (p *Player) AddEntityOnMap(e EntityOnMap) {
	p.EntitiesOnSameMap = append(p.EntitiesOnSameMap, e)
}

func (p *Player) RemoveEntityOnMap(id int) {
	var nouv []EntityOnMap
	for _, e := range p.EntitiesOnSameMap {
		if e.Id != id {
			nouv = append(nouv, e)
		}
	}
	p.EntitiesOnSameMap = nouv
}

func (p *Player) ClearEntityOnMap() {
	p.EntitiesOnSameMap = []EntityOnMap{}
}

func (p *Player) UpdateCellId(idEntity int, cell int) {
	for i, e := range p.EntitiesOnSameMap {
		if e.Id == idEntity {
			fmt.Println("UpdateCellId of entity id", idEntity, "cell id", cell)
			p.EntitiesOnSameMap[i].CellId = cell
			return
		}
	}
	fmt.Println("Entity", idEntity, "not found")
	return
}

func (e EntityOnMap) GetLevel() int {
	sum := 0
	for _, lvl := range e.levels {
		sum += lvl
	}
	return sum
}

func (p Player) GetAFigthableEntity() (EntityOnMap, error) {
	if len(p.EntitiesOnSameMap) == 0 {
		return EntityOnMap{}, errors.New("No entity on this map")
	} else {
		for _, e := range p.EntitiesOnSameMap {
			if e.Id < 0 && e.GetLevel() > 1 { // Avoid PNJ
				return e, nil
			}
		}
	}
	return EntityOnMap{}, errors.New("No entity found on the map")
}

func BuildEntity(datas []string) EntityOnMap {

	cellId, _ := strconv.Atoi(datas[0][1:])
	Id, _ := strconv.Atoi(datas[3])
	levelsStr := strings.Split(datas[7], ",")

	var levelsInt []int
	for _, level := range levelsStr {
		levelInt, _ := strconv.Atoi(level)
		levelsInt = append(levelsInt, levelInt)
	}

	return EntityOnMap{
		CellId: cellId,
		Id:     Id,
		levels: levelsInt,
	}
}