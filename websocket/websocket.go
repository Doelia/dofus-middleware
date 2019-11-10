package websocket

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"sync"
)


type WebSocket struct {

	OnConnexion func()

	OnMessage func(message string)

	mutex *sync.Mutex

	conn *websocket.Conn

}

var upgrader = websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true} } // use default options


func (s *WebSocket) sendPacket(typepacket string, content string) {
	if s == nil || s.conn == nil {
		return
	}

	message := []byte(typepacket + "|" + content)
	//fmt.Println("web.SendMessage", string(message))

	s.mutex.Lock()
	err := s.conn.WriteMessage(websocket.TextMessage, message)
	if err != nil {
		log.Println("write error: ", err)
	}
	s.mutex.Unlock()
}


func (s *WebSocket) echo(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	s.conn = c
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()
	s.OnConnexion()
	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		if mt == websocket.TextMessage {
			s.OnMessage(string(message))
		}
	}
}

func (s *WebSocket) StartWebSocket() {

	if s.OnConnexion == nil || s.OnMessage == nil {
		fmt.Println("Cant startWebSocket, OnConnexion on OnMessage handler is not defined")
		return;
	}

	s.mutex = &sync.Mutex{}

	fmt.Println("Start web socket server")
	http.HandleFunc("/ws", s.echo)
	log.Fatal(http.ListenAndServe("localhost:8001", nil))
}