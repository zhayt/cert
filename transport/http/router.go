package http

import (
	"github.com/gorilla/mux"
	"net/http"
)

func (s *Server) InitRoute() http.Handler {
	r := mux.NewRouter()

	r.HandleFunc("/rest/substr/find", s.handler.FindSubString).Methods(http.MethodPost)

	r.HandleFunc("/rest/email/check", s.handler.AnalysisToEmail).Methods(http.MethodPost)

	r.HandleFunc("/rest/counter/add/{i}", s.handler.CounterIncrease).Methods(http.MethodPost)
	r.HandleFunc("/rest/counter/sub/{i}", s.handler.CounterDecrease).Methods(http.MethodPost)
	r.HandleFunc("/rest/counter/val", s.handler.ShowCounter).Methods(http.MethodGet)

	r.HandleFunc("/rest/user", s.handler.CreateUser).Methods(http.MethodPost)
	r.HandleFunc("/rest/user/{id}", s.handler.GetUser).Methods(http.MethodGet)
	r.HandleFunc("/rest/user/{id}", s.handler.UpdateUser).Methods(http.MethodPut)
	r.HandleFunc("/rest/user/{id}", s.handler.DeleteUser).Methods(http.MethodDelete)

	r.HandleFunc("/rest/hash/calc", s.handler.HashCalculation).Methods(http.MethodPost)
	r.HandleFunc("/rest/hash/result/{id}", s.handler.ShowHash).Methods(http.MethodGet)

	return r
}
