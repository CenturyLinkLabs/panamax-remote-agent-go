package main // import "github.com/CenturyLinkLabs/panamax-remote-agent-go"

import (
	"flag"
	"log"
	"os"
	"regexp"

	"github.com/CenturyLinkLabs/panamax-remote-agent-go/adapter"
	"github.com/CenturyLinkLabs/panamax-remote-agent-go/agent"
	"github.com/CenturyLinkLabs/panamax-remote-agent-go/api"
	"github.com/CenturyLinkLabs/panamax-remote-agent-go/repo"
)

func main() {
	sf := flag.Bool("secure", true, "whether or not to use SSL and BasicAuth, defaults to true")
	flag.Parse()

	p, err := repo.MakeDeploymentStore(dbLocation())
	if err != nil {
		log.Fatal(err)
	}
	c := adapter.MakeClient(adapterEndpoint())
	dm := agent.MakeDeploymentManager(p, c, version())
	s := makeServer(dm, sf)
	s.Start(serverPort())
}

func makeServer(dm agent.Manager, sf *bool) api.Server {
	log.Printf("secure?: %t", *sf)

	if *sf {
		return api.MakeServer(
			dm,
			username(),
			password(),
			certFile(),
			keyFile(),
		)
	} else {
		return api.MakeInsecureServer(dm)
	}
}

func version() string {
	return os.Getenv("REMOTE_AGENT_VERSION")
}

func username() string {
	return os.Getenv("REMOTE_AGENT_ID")
}

func password() string {
	return os.Getenv("REMOTE_AGENT_API_KEY")
}

func keyFile() string {
	return "/usr/local/share/certs/pmx_remote_agent.key"
}

func certFile() string {
	return "/usr/local/share/certs/pmx_remote_agent.crt"
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
