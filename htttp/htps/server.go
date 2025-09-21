package htttp

import (
	"net/http"

	"github.com/gorilla/mux"
)

type HTTPServer struct {
	HTTPHandler HTTPHandlers
}

func NewHTTPServer(HTTPHandlers HTTPHandlers) HTTPServer {
	return HTTPServer{
		HTTPHandler: HTTPHandlers,
	}
}

func (s *HTTPServer) Start() {
	router := mux.NewRouter()

	router.Path("/users").Methods(http.MethodGet).HandlerFunc(s.HTTPHandler.HandlerGetAllUsers)
	router.Path("/users").Methods(http.MethodPost).HandlerFunc(s.HTTPHandler.HandlerNewUser)
	router.Path("/users/cash/{username}").Methods(http.MethodPatch).HandlerFunc(s.HTTPHandler.HandlerCashReciver)
	router.Path("/users/transfer/{username}").Methods(http.MethodPatch).HandlerFunc(s.HTTPHandler.HandlerTransferOperation)
	router.Path("/users/buy/{username}").Methods(http.MethodPatch).HandlerFunc(s.HTTPHandler.HandlerBuynigOperation)
	router.Path("/users/{username}").Methods(http.MethodDelete).HandlerFunc(s.HTTPHandler.HandlerDeleteUser)

	http.ListenAndServe(":8080", router)

}
