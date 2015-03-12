package agent

type Deployment struct {
	Override Template `json:"override,omitempty"`
	Template Template `json:"template,omitempty"`
}

func (d *Deployment) MergedImages() []Image {
	mImgs := make([]Image, len(d.Template.Images))

	for i, tImg := range d.Template.Images {
		for _, oImg := range d.Override.Images {
			if oImg.Name == tImg.Name {
				tImg.OverrideWith(oImg)
			}
		}
		mImgs[i] = tImg
	}
	return mImgs
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

func (img *Image) OverrideWith(o Image) {
	img.overrideSource(o)
	img.overrideEnv(o)
	// img.overrideDeployment(o)
	// img.Deployment = o.Deployment
	// img.Links = o.Links
	// img.Ports = o.Ports
}

func (i *Image) overrideDeployment(o Image) {

}

func (img *Image) overrideEnv(o Image) {
	//TODO add the extra override envs that didn't exist in base
	envs := make([]Environment, 0)

	for _, env := range img.Environment {
		for _, oEnv := range o.Environment {
			if env.Variable == oEnv.Variable {
				env = oEnv
			}
		}
		envs = append(envs, env)
	}
	img.Environment = envs
}

func (i *Image) overrideSource(o Image) {
	if o.Source != "" {
		i.Source = o.Source
	}
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
