package main

import (
	// "os"
	"log"
	"regexp"

	"github.com/CenturyLinkLabs/panamax-remote-agent-go/adapter"
	"github.com/CenturyLinkLabs/panamax-remote-agent-go/agent"
	"github.com/CenturyLinkLabs/panamax-remote-agent-go/api"
	"github.com/CenturyLinkLabs/panamax-remote-agent-go/repo"
)

func main() {
	p, err := repo.MakePersister("db/agent.db")
	if err != nil {
		log.Fatal(err)
	}
	c := adapter.MakeClient(adapterEndpoint())
	dm := agent.MakeDeploymentManager(p, c)
	s := api.NewServer(dm)
	s.Start(":1234")
}

func adapterEndpoint() string {
	// url := os.Getenv("ADAPTER_PORT")
	url := "tcp://localhost:1234"
	r, _ := regexp.Compile("^tcp")
	return r.ReplaceAllString(url, "http")
}
