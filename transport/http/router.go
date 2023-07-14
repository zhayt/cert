package http

import (
	"github.com/gorilla/mux"
	"net/http"
)

func (s *Server) InitRoute() http.Handler {
	r := mux.NewRouter()

	r.HandleFunc("/rest/user", s.handler.CreateUser).Methods(http.MethodPost)
	r.HandleFunc("/rest/user/{id}", s.handler.GetUser).Methods(http.MethodGet)
	r.HandleFunc("/rest/user/{id}", s.handler.UpdateUser).Methods(http.MethodPut)
	r.HandleFunc("/rest/user/{id}", s.handler.DeleteUser).Methods(http.MethodDelete)

	return r
}
