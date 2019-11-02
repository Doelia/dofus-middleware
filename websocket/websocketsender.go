package websocket

import (
	"dofusmiddleware/options"
	"dofusmiddleware/world"
	"encoding/json"
)

func (s *WebSocket) SendCharacters(characters []world.Player) {
	message, _ := json.Marshal(characters)
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
