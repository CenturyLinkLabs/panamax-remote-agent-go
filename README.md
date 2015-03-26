# panamax-remote-agent-go

## Installation for development

Pre-requisites
* Golang development environment.
* SQLite3 (https://sqlite.org/)

```bash
$ go get github.com/CenturyLinkLabs/panamax-remote-agent-go
$ cd $GOPATH/src/github.com/CenturyLinkLabs/panamax-remote-agent-go
$ sqlite3 db/agent.db < db/create.sql # for main
$ sqlite3 db/agent_test.db < db/create.sql # for integration tests
$ go test ./... #should all pass

```
