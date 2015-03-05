package api

import "net/http"
import "github.com/CenturyLinkLabs/panamax-remote-agent-go/agent"

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc func(agent.DeploymentManager, http.ResponseWriter, *http.Request)
}

type Routes []Route

var routes = Routes{
	{
		"ShowDeployment",
		"GET",
		"/deployments/{id}",
		ShowDeployment,
	},
	{
		"ListDeployments",
		"GET",
		"/deployments",
		ListDeployments,
	},
	{
		"CreateDeployment",
		"POST",
		"/deployments",
		CreateDeployment,
	},
	{
		"DeleteDeployment",
		"DELETE",
		"/deployments/{id}",
		DeleteDeployment,
	},
}
