package main

import (
	"context"
	"flag"
	"github.com/OrlovEvgeny/sctun/server"
	"log"
)

const defaultAddr = "0.0.0.0:8080"

var addr string

func init() {
	flag.StringVar(&addr, "addr", defaultAddr, "--addr "+defaultAddr)
	flag.Parse()
}

func main() {
	log.Println(addr)
	srv := server.Server{}
	if err := srv.Run(context.Background(), addr); err != nil {
		log.Fatal(err)
	}
}
