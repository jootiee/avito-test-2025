package httpserver

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type server struct {
	router *mux.Router
	logger *logrus.Logger
	// storage
}

func newServer() *server {
	s := &server{
		router: mux.NewRouter(),
		logger: logrus.New(),
	}

	s.configureRouter()

	return s
}

func (s *server) configureRouter(){
	// set request id...
	// log request...
	s.router.HandleFunc("/", s.handleHello()).Methods("GET")
}
func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *server) handleHello() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		type response struct {
			Message string `json:"message"`
		}
		resp := &response{Message: "Hello, World!"}
		s.respond(w, r, http.StatusOK, resp)	
	}
}
	

func (s * server) respond(w http.ResponseWriter, r *http.Request, code int, data interface{}) {
	w.WriteHeader(code)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}