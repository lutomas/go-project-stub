package server

import (
	"net/http"

	"github.com/lutomas/go-project-stub/types"
)

func (s *Server) appVersionHandler(resp http.ResponseWriter, req *http.Request) {
	response(types.NewVersion("MAIN-APP"), http.StatusOK, nil, resp, req)
}
