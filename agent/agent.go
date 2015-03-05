package agent

type Deployment struct {
	Override template `json:"override,omitempty"`
	Template template `json:"template,omitempty"`
}

type DeploymentResponses []DeploymentResponseLite

type DeploymentResponseLite struct {
	ID           int      `json:"id"`
	Name         string   `json:"name"`
	Redeployable bool     `json:"redeployable"`
	ServiceIDs   []string `json:"service_ids"`
}

type DeploymentResponseFull struct {
	ID           int           `json:"id"`
	Name         string        `json:"name"`
	Redeployable bool          `json:"redeployable"`
	Status       serviceStatus `json:"status"`
}

type environment struct {
	Variable string `json:"variable,omitempty"`
	Value    string `json:"value,omitempty"`
}

type link struct {
	Service string `json:"service,omitempty"`
	Alias   string `json:"alias,omitempty"`
}

type port struct {
	HostPort      int `json:"host_port,omitempty"`
	ContainerPort int `json:"container_port,omitempty"`
}

type deploymentSettings struct {
	Count int `json:"count"`
}

type Image struct {
	Name        string             `json:"name,omitempty"`
	Source      string             `json:"source,omitempty"`
	Deployment  deploymentSettings `json:"deployment,omitempty"`
	Links       []link             `json:"links,omitempty"`
	Environment []environment      `json:"environment,omitempty"`
	Ports       []port             `json:"port,omitemptys`
}

type template struct {
	Name   string  `json:"name,omitempty"`
	Images []Image `json:"images,omitempty"`
	// TODO: Description?
}

type service struct {
	ID          string `json:"id"`
	ActualState string `json:"actualState"`
}

type services []service

type serviceStatus struct {
	Services services `json:"services"`
}
