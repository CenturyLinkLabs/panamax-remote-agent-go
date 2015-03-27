package api

import (
	"log"
	"net/http"
	"time"

	"github.com/CenturyLinkLabs/panamax-remote-agent-go/agent"
	"github.com/gorilla/mux"
)

// A Server is the HTTP server which responds to API requests.
type Server struct {
	DeploymentManager agent.DeploymentManager
	username          string
	password          string
}

// MakeServer returns a new Server instance containting a manager to which it will defer work.
func MakeServer(dm agent.DeploymentManager, un string, pw string) Server {
	return Server{
		DeploymentManager: dm,
		username:          un,
		password:          pw,
	}
}

func (s Server) newRouter() *mux.Router {
	r := mux.NewRouter()

	dm := s.DeploymentManager

	for _, route := range routes {
		fct := route.HandlerFunc
		wrap := func(w http.ResponseWriter, r *http.Request) {

			if !s.isAuthenticated(r) {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			// make it json
			w.Header().Set("Content-Type", "application/json; charset=utf-8")

			// log it
			st := time.Now()

			log.Printf(
				"%s\t%s\t%s\t%s",
				r.Method,
				r.RequestURI,
				route.Name,
				time.Since(st),
			)

			fct(dm, w, r)
		}

		r.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			HandlerFunc(wrap)
	}

	return r
}

// Start initializes a router and starts the server at a given address.
func (s Server) Start(addr string) {
	r := s.newRouter()
	log.Fatal(http.ListenAndServe(addr, r))
}

func (s Server) isAuthenticated(r *http.Request) bool {
	un, pw, ok := r.BasicAuth()

	if ok && (un == s.username) && (pw == s.password) {
		return true
	}

	return false
}
