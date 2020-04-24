package server

import (
	"context"
	"log"
	"net/http"
)

//Run http server start
func (s *Server) RunHttp(ctx context.Context, addr string) {
	http.HandleFunc("/", s.getServers)

	log.Printf("http server start on %s\n", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("Error in ListenAndServe: %s", err)
	}
}

func (s *Server) getServers(w http.ResponseWriter, r *http.Request) {

}
