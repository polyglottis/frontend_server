// Package server contains helper functions and structs for the frontend server.
package server

import (
	"bytes"
	"io"
	"log"

	"github.com/polyglottis/platform/frontend"
)

func Call(context *frontend.Context, f func(io.Writer, *TmplArgs) error) (answer []byte, err error) {
	args, err := GetTmplArgs(context)
	if err != nil {
		return nil, err
	}

	buffer := new(bytes.Buffer)
	err = f(buffer, args)
	if err != nil {
		log.Println("Error:", err)
		return nil, err
	}
	return buffer.Bytes(), nil
}
