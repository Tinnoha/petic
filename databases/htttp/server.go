package htttp

import (
	"net/http"

	"github.com/gorilla/mux"
)

type HttpServer struct {
	HttpHandler HttpHandler
}

func NewHttpServer(Hh HttpHandler) *HttpServer {
	return &HttpServer{
		HttpHandler: Hh,
	}
}

func (s *HttpServer) Run() {
	router := mux.NewRouter()

	router.Path("/users").Methods(http.MethodGet).HandlerFunc(s.HttpHandler.GetUsers)
	router.Path("/users").Methods(http.MethodPost).HandlerFunc(s.HttpHandler.Insert)

	router.Path("/users/cash/{Username}").Methods(http.MethodPatch).HandlerFunc(s.HttpHandler.AddCash)
	router.Path("/users/buy/{Username}").Methods(http.MethodPatch).HandlerFunc(s.HttpHandler.Buy)
	router.Path("/users/transfer/{Username}").Methods(http.MethodPatch).HandlerFunc(s.HttpHandler.TransferCash)

	router.Path("/users/{Username}").Methods(http.MethodDelete).HandlerFunc(s.HttpHandler.DeleteUser)

	http.ListenAndServe(":8081", router)
}
