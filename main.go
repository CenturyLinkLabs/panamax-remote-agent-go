package main

import (
	// "os"
	"regexp"

	"github.com/CenturyLinkLabs/panamax-remote-agent-go/adapter"
	"github.com/CenturyLinkLabs/panamax-remote-agent-go/agent"
	"github.com/CenturyLinkLabs/panamax-remote-agent-go/api"
	"github.com/CenturyLinkLabs/panamax-remote-agent-go/repo"
)

func main() {
	dRepo := repo.NewDeploymentRepo("db/agent.db")
	ad := adapter.MakeClient(adapterEndpoint())
	dm := agent.MakeDeploymentManager(dRepo, ad)
	s := api.NewServer(dm)
	s.Start(":1234")
}

func adapterEndpoint() string {
	// url := os.Getenv("ADAPTER_PORT")
	url := "tcp://localhost:1234"
	r, _ := regexp.Compile("^tcp")
	return r.ReplaceAllString(url, "http")
}
