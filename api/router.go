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
}

// MakeServer returns a new Server instance containting a manager to which it will defer work.
func MakeServer(dm agent.DeploymentManager) Server {
	return Server{DeploymentManager: dm}
}

func (s Server) newRouter() *mux.Router {
	r := mux.NewRouter()

	dm := s.DeploymentManager

	for _, route := range routes {
		fct := route.HandlerFunc
		wrap := func(w http.ResponseWriter, r *http.Request) {
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
