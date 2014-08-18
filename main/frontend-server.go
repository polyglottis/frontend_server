// Package main contains the frontend-server executable.
package main

import (
	"log"

	"github.com/polyglottis/frontend_server/register"
	"github.com/polyglottis/platform/config"
	"github.com/polyglottis/platform/frontend/rpc"
)

func main() {
	c := config.Get()

	main := register.NewServer()
	s := rpc.NewFrontendServer(main, c.Frontend)

	err := s.RegisterAndListen()
	if err != nil {
		log.Fatalln(err)
	}
	defer s.Close()
	log.Printf("Frontend Server listening on %v", c.Frontend)

	s.Accept()
}
