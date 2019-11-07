package websocket

import (
	"dofusmiddleware/options"
	"dofusmiddleware/world"
	"encoding/json"
)

func (s *WebSocket) SendCharacters(characters []*world.Player) {

	// remove unwanted fields
	var toSend []world.Player
	for _, p := range characters {
		p.Connexion = nil
		toSend = append(toSend, *p)
	}

	message, _ := json.Marshal(toSend)
	s.sendPacket("CHARACTERS", string(message))
}

func (s *WebSocket) SendOptions(options options.OptionsStruct) {
	message, _ := json.Marshal(options)
	s.sendPacket("OPTIONS", string(message))
}

func (s *WebSocket) SendMap(themap world.Map) {
	message, _ := json.Marshal(themap)
	s.sendPacket("MAP", string(message))
}

func (s *WebSocket) SendPath(path []int) {
	message, _ := json.Marshal(path)
	s.sendPacket("PATH", string(message))
}
