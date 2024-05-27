package server

import (
	"errors"
	"net/http"
)

var (
	ErrUnexpectedDstUUIDOnReturn = errors.New("unexpected destination uuid on probe's return")
	ErrUnexpectedSrcDstUUIDs     = errors.New("source uuid is not us, but non-zero destination uuid")
)

func (s *Server) handleHealthcheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
