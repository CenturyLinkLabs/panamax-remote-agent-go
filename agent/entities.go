package agent

import (
	"encoding/json"
)

type DeploymentBlueprint struct {
	Override Template `json:"override,omitempty"`
	Template Template `json:"template,omitempty"`
}

func (d *DeploymentBlueprint) MergedImages() []Image {
	mImgs := make([]Image, len(d.Template.Images), len(d.Template.Images))

	for i, tImg := range d.Template.Images {
		for _, oImg := range d.Override.Images {
			if oImg.Name == tImg.Name {
				tImg.overrideWith(oImg)
			}
		}

		mImgs[i] = tImg
	}
	return mImgs
}

type Template struct {
	Name   string  `json:"name,omitempty"`
	Images []Image `json:"images,omitempty"`
}

type Image struct {
	Name        string
	Source      string
	Command     string
	Deployment  DeploymentSettings
	Links       []Link
	Environment []Environment
	Ports       []Port
	Expose      []int
	Volumes     []Volume
	VolumesFrom []string
}

func (img Image) MarshalJSON() ([]byte, error) {
	i := map[string]interface{}{}

	if img.Name != "" {
		i["name"] = img.Name
	}
	if img.Source != "" {
		i["source"] = img.Source
	}
	if img.Command != "" {
		i["command"] = img.Command
	}
	if (img.Deployment != DeploymentSettings{}) {
		i["deployment"] = img.Deployment
	}
	if len(img.Links) > 0 {
		i["links"] = img.Links
	}
	if len(img.Environment) > 0 {
		i["environment"] = img.Environment
	}
	if len(img.Ports) > 0 {
		i["ports"] = img.Ports
	}
	if len(img.Expose) > 0 {
		i["expose"] = img.Expose
	}
	if len(img.Volumes) > 0 {
		i["volumes"] = img.Volumes
	}
	if len(img.VolumesFrom) > 0 {
		i["volumesFrom"] = img.VolumesFrom
	}

	return json.Marshal(i)
}

func (img *Image) overrideWith(o Image) {
	img.overrideEnv(o)
	img.overrideDeployment(o)
}

func (img *Image) overrideDeployment(o Image) {
	if (o.Deployment != DeploymentSettings{}) {
		img.Deployment = o.Deployment
	}
}

func (img *Image) overrideEnv(o Image) {
	var envs []Environment

	for _, env := range img.Environment {
		for i, oEnv := range o.Environment {
			if env.Variable == oEnv.Variable {
				env = oEnv
				o.Environment = append(o.Environment[:i], o.Environment[i+1:]...)
			}
		}
		envs = append(envs, env)
		envs = append(envs, o.Environment...)
	}
	img.Environment = envs
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
	HostPort      int `json:"hostPort,omitempty"`
	ContainerPort int `json:"containerPort,omitempty"`
}

type DeploymentSettings struct {
	Count int `json:"count,omitempty"`
}

type Volume struct {
	ContainerPath string `json:"containerPath"`
	HostPath      string `json:"hostPath"`
}

type DeploymentResponseLite struct {
	ID           int      `json:"id"`
	Name         string   `json:"name"`
	Redeployable bool     `json:"redeployable"`
	ServiceIDs   []string `json:"service_ids"`
}

type DeploymentResponseFull struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	Redeployable bool   `json:"redeployable"`
	Status       Status `json:"status"`
}

type Status struct {
	Services []Service `json:"services"`
}

type Service struct {
	ID          string `json:"id"`
	ActualState string `json:"actualState"`
}

type Metadata struct {
	Agent struct {
		Version string `json:"version"`
	} `json:"agent"`
	Adapter interface{} `json:"adapter"`
}
