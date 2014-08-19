// Package register contains the server definition (func register).
package register

import (
	"github.com/polyglottis/frontend_server"
	"github.com/polyglottis/frontend_server/extract"
	"github.com/polyglottis/frontend_server/extract/edit"
	"github.com/polyglottis/frontend_server/user"
)

type Server struct {
	*frontend_server.Server
	*extract.ExtractServer
	*edit.EditServer
	*user.UserServer
}

func NewServer() *Server {
	return &Server{
		Server:        &frontend_server.Server{},
		ExtractServer: &extract.ExtractServer{},
		EditServer:    &edit.EditServer{},
		UserServer:    &user.UserServer{},
	}
}
