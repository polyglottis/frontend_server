// Package main contains the frontend-server executable.
package main

import (
	"flag"
	"log"

	"github.com/polyglottis/frontend_server"
	"github.com/polyglottis/platform/frontend/rpc"
)

var tcpAddr = flag.String("tcp", ":18658", "TCP address of frontend server")

func main() {
	flag.Parse()

	main := frontend_server.New()
	s := rpc.NewFrontendServer(main, *tcpAddr)

	err := s.RegisterAndListen()
	if err != nil {
		log.Fatalln(err)
	}
	defer s.Close()
	log.Printf("Frontend Server listening on %v", *tcpAddr)

	s.Accept()
}
