package server

import (
	"context"
	"fmt"
	"github.com/OrlovEvgeny/sctun/netpack"
	"log"
	"net"
)

type Server struct {
	conns      chan net.Conn
	errs       chan error
	hub        *netpack.Hub
	addr       string
	ExternalIP string
}

//Run start master tunnel
func (s *Server) Run(ctx context.Context, addr string) (err error) {
	ctx, cancel := context.WithCancel(ctx)
	s.conns = make(chan net.Conn, 100)
	s.errs = make(chan error, 100)
	s.hub = netpack.NewHub()
	s.hub.ExternalIP = s.ExternalIP
	s.addr = addr

	go s.listen()

	defer func() {
		cancel()
		close(s.conns)
	}()

	for {
		select {
		case <-ctx.Done():
			err = ctx.Err()
			return
		case conn := <-s.conns:
			s.hub.JoinToSrv(conn)
		}
	}
}

//listen tcp master server
func (s *Server) listen() {
	addr, _ := net.ResolveTCPAddr("tcp", s.addr)
	log.Printf("starting master-server on %s\n", addr.String())
	listener, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return
	}
	defer listener.Close()

	for {
		c, err := listener.Accept()
		if err != nil {
			s.errs <- fmt.Errorf("error accepting connection %v", err)
			continue
		}

		log.Printf("accepted new slave proxy node connection from %v", c.RemoteAddr())
		s.conns <- c
	}
}
