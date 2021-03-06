[![](https://badge.imagelayers.io/centurylink/panamax-remote-agent-go.svg)](https://imagelayers.io/?images=centurylink/panamax-remote-agent-go:latest 'Get your own badge on imagelayers.io')

![Panamax - Docker Management for Humans](http://panamax.ca.tier3.io/panamax_ui_wiki_screens/panamax_logo-title.png)

[Panamax](http://panamax.io) is a containerized app creator with an open-source app marketplace hosted in GitHub. Panamax provides a friendly interface for users of Docker, Fleet & CoreOS. With Panamax, you can easily create, share, and deploy any containerized app no matter how complex it might be. Learn more at [Panamax.io](http://panamax.io) or browse the [Panamax Wiki](https://github.com/CenturyLinkLabs/panamax-ui/wiki).

# panamax-remote-agent-go

[![Circle CI](https://circleci.com/gh/CenturyLinkLabs/panamax-remote-agent-go.svg?style=svg)](https://circleci.com/gh/CenturyLinkLabs/panamax-remote-agent-go)

Replaces [CenturyLinkLabs/panamax-remote-agent](https://github.com/CenturyLinkLabs/panamax-remote-agent)

## NOTE

This repo is no longer being maintained. Users are welcome to fork it, but we make no warranty of its functionality.

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
