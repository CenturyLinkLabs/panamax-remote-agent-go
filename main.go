package main

import (
	"log"
	// "os"
	"regexp"

	"github.com/CenturyLinkLabs/panamax-remote-agent-go/adapter"
	"github.com/CenturyLinkLabs/panamax-remote-agent-go/agent"
	"github.com/CenturyLinkLabs/panamax-remote-agent-go/api"
	"github.com/CenturyLinkLabs/panamax-remote-agent-go/repo"
)

func main() {
	dRepo := repo.NewDeploymentRepo("db/agent.db")

	ad := adapter.NewClient(adapterEndpoint())

	dm, err := agent.NewDeploymentManager(dRepo, ad)
	if err != nil {
		log.Fatal(err)
	}

	s := api.NewServer(dm)
	s.Start(":1234")
}

func adapterEndpoint() string {
	// url := os.Getenv("ADAPTER_PORT")
	url := "tcp://localhost:1234"
	r, _ := regexp.Compile("^tcp")
	return r.ReplaceAllString(url, "http")
}
