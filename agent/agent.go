package agent

type Deployment struct {
	Override Template `json:"override,omitempty"`
	Template Template `json:"template,omitempty"`
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
	Status       ServiceStatus `json:"status"`
}

type Environment struct {
	Variable string `json:"variable,omitempty"`
	Value    string `json:"value,omitempty"`
}

type Link struct {
	Service string `json:"service,omitempty"`
	Alias   string `json:"alias,omitempty"`
}

type Port struct {
	HostPort      int `json:"host_port,omitempty"`
	ContainerPort int `json:"container_port,omitempty"`
}

type DeploymentSettings struct {
	Count int `json:"count"`
}

type Image struct {
	Name        string             `json:"name,omitempty"`
	Source      string             `json:"source,omitempty"`
	Deployment  DeploymentSettings `json:"deployment,omitempty"`
	Links       []Link             `json:"links,omitempty"`
	Environment []Environment      `json:"environment,omitempty"`
	Ports       []Port             `json:"port,omitemptys`
}

type Template struct {
	Name   string  `json:"name,omitempty"`
	Images []Image `json:"images,omitempty"`
	// TODO: Description?
}

type service struct {
	ID          string `json:"id"`
	ActualState string `json:"actualState"`
}

type Services []service

type ServiceStatus struct {
	Services Services `json:"services"`
}
