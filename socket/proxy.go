package socket

import (
	"crypto/tls"
	"dofusmiddleware/world"
	"fmt"
	"log"
	"net"
)


// Server is a TCP WebSocket that takes an incoming request and sends it to another
// WebSocket, proxying the response back to the client.
type Server struct {
	// TCP address to listen on
	Addr string

	// TCP address of target WebSocket
	Target string

	// ModifyRequest is an optional function that modifies the request from a client to the target WebSocket.
	ModifyRequest func(b *[]byte, id string)

	// ModifyResponse is an optional function that modifies the response from the target WebSocket.
	ModifyResponse func(b *[]byte, id string)

	// TLS configuration to listen on.
	TLSConfig *tls.Config

	// TLS configuration for the proxy if needed to connect to the target WebSocket with TLS protocol.
	// If nil, TCP protocol is used.
	TLSConfigTarget *tls.Config
}

// ListenAndServe listens on the TCP network address laddr and then handle packets
// on incoming connections.
func (s *Server) ListenAndServe() error {
	listener, err := net.Listen("tcp", s.Addr)
	if err != nil {
		fmt.Println("Error")
		return err
	}
	return s.serve(listener)
}

func (s *Server) serve(ln net.Listener) error {
	for {
		conn, err := ln.Accept()
		fmt.Println("serve")
		if err != nil {
			log.Println(err)
			continue
		}
		go s.handleConn(conn)
	}
}


func (s *Server) handleConn(conn net.Conn) {

	fmt.Println("handleConn")
	// connects to target WebSocket
	var rconn net.Conn
	var err error

	rconn, err = net.Dial("tcp", s.Target)
	if err != nil {
		fmt.Println(err)
		return
	}

	// write to dst what it reads from src
	var pipe = func(src, dst net.Conn, filter func(b *[]byte, id string), sens string) {
		defer func() {
			_ = conn.Close()
			_ = rconn.Close()
		}()

		buff := make([]byte, 65535)
		for {
			n, err := src.Read(buff)

			if err != nil {
				log.Println(err)
				return
			}
			b := buff[:n]

			var id string
			if sens == "in" {
				id = src.RemoteAddr().String()
			} else {
				id = dst.RemoteAddr().String()
			}

			if filter != nil {
				filter(&b, id)
			}

			_, err = dst.Write(b)
			if err != nil {
				log.Println(err)
				return
			}
		}
	}

	fmt.Println("Add character")
	world.Characters = append(world.Characters, world.Character{
		Id: conn.RemoteAddr().String(),
		ConnClient: conn,
		ConnServer: rconn,
	})

	go pipe(conn, rconn, s.ModifyRequest, "in")
	go pipe(rconn, conn, s.ModifyResponse, "out")

	fmt.Println("pipe started")
}