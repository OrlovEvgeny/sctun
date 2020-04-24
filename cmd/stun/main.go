package main

import (
	"context"
	"flag"
	"github.com/OrlovEvgeny/sctun/server"
	"log"
)

const defaultAddr = "0.0.0.0:8080"

var (
	addr     string
	external string
	httpAddr string
)

func init() {
	flag.StringVar(&addr, "addr", defaultAddr, "-addr <ip_addr>")
	flag.StringVar(&external, "external", "127.0.0.1", "-external <ip_addr>")
	flag.StringVar(&httpAddr, "http", "0.0.0.0:8181", "-http <ip_addr>")
	flag.Parse()

	log.SetFlags(log.Lmicroseconds | log.Lshortfile)
}

func main() {
	srv := &server.Server{ExternalIP: external}
	go func(srv *server.Server) {
		srv.RunHttp(httpAddr)
	}(srv)
	if err := srv.Run(context.Background(), addr); err != nil {
		log.Fatal(err)
	}
}
