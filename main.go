package main // import "github.com/CenturyLinkLabs/panamax-remote-agent-go"

import (
	"log"
	"os"
	"regexp"

	"github.com/CenturyLinkLabs/panamax-remote-agent-go/adapter"
	"github.com/CenturyLinkLabs/panamax-remote-agent-go/agent"
	"github.com/CenturyLinkLabs/panamax-remote-agent-go/api"
	"github.com/CenturyLinkLabs/panamax-remote-agent-go/repo"
)

func main() {
	p, err := repo.MakePersister(dbLocation())
	if err != nil {
		log.Fatal(err)
	}
	c := adapter.MakeClient(adapterEndpoint())
	dm := agent.MakeDeploymentManager(p, c)
	s := api.MakeServer(dm)
	s.Start(serverPort())
}

func dbLocation() string {
	l := os.Getenv("DB_LOCATION")
	if l == "" {
		l = "db/agent.db"
	}
	return l
}

func serverPort() string {
	p := os.Getenv("SERVER_PORT")
	if p == "" {
		p = "3000"
	}
	return ":" + p
}

func adapterEndpoint() string {
	url := os.Getenv("ADAPTER_PORT")
	if url == "" {
		url = "tcp://localhost:3001"
	}
	r, _ := regexp.Compile("^tcp")
	return r.ReplaceAllString(url, "http")
}
