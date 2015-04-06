package api

import "net/http"
import "github.com/CenturyLinkLabs/panamax-remote-agent-go/agent"

type route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc func(agent.Manager, http.ResponseWriter, *http.Request)
}

var routes = []route{
	{
		"showDeployment",
		"GET",
		"/deployments/{id}",
		showDeployment,
	},
	{
		"listDeployments",
		"GET",
		"/deployments",
		listDeployments,
	},
	{
		"createDeployment",
		"POST",
		"/deployments",
		createDeployment,
	},
	{
		"deleteDeployment",
		"DELETE",
		"/deployments/{id}",
		deleteDeployment,
	},
	{
		"reDeploy",
		"POST",
		"/deployments/{id}/redeploy",
		reDeploy,
	},
	{
		"metadata",
		"GET",
		"/metadata",
		metadata,
	},
}
