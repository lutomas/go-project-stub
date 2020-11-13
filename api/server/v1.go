package server

import "github.com/gorilla/mux"

func (s *Server) infoHandlers(mux *mux.Router) {
	// application version handler
	mux.HandleFunc("/v1/info/version", s.appVersionHandler).Methods("GET")
}
