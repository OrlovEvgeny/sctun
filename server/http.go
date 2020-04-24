package server

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

//Run http server start
func (s *Server) RunHttp(addr string) {
	http.HandleFunc("/", s.getServers)

	log.Printf("http server start on %s\n", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("Error in ListenAndServe: %s", err)
	}
}

//getServers
func (s *Server) getServers(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, strings.Join(s.hub.ProxyList(), "\n"))
}
