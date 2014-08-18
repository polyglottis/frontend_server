// Package register contains the server definition (func register).
package register

import (
	"github.com/polyglottis/frontend_server"
	"github.com/polyglottis/frontend_server/extract"
	"github.com/polyglottis/frontend_server/extract/edit"
)

type Server struct {
	*frontend_server.Server
	*extract.ExtractServer
	*edit.EditServer
}

func NewServer() *Server {
	return &Server{
		Server:        &frontend_server.Server{},
		ExtractServer: &extract.ExtractServer{},
		EditServer:    &edit.EditServer{},
	}
}
