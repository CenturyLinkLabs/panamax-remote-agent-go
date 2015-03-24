package agent

import (
	"encoding/json"
)

type DeploymentBlueprint struct {
	Override Template `json:"override,omitempty"`
	Template Template `json:"template,omitempty"`
}

func (d *DeploymentBlueprint) MergedImages() []Image {
	mImgs := make([]Image, len(d.Template.Images))

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
	Name        string             `json:"name,omitempty"`
	Source      string             `json:"source,omitempty"`
	Command     string             `json:"command,omitempty"`
	Deployment  DeploymentSettings `json:"deployment,omitempty"`
	Links       []Link             `json:"links,omitempty"`
	Environment []Environment      `json:"environment,omitempty"`
	Ports       []Port             `json:"port,omitemptys`
	Expose      []int              `json:"expose"`
	Volumes     []Volume           `json:"volumes"`
	VolumesFrom []string           `json:"volumesFrom"`
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
	//TODO add the extra override envs that didn't exist in base
	var envs []Environment

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

type DeploymentResponses []DeploymentResponseLite

func NewDeploymentResponseLite(id int, nm string, tpl string, sids string) *DeploymentResponseLite {
	drl := &DeploymentResponseLite{
		ID:           id,
		Name:         nm,
		Redeployable: tpl != "",
	}

	json.Unmarshal([]byte(sids), &drl.ServiceIDs)
	json.Unmarshal([]byte(tpl), &drl.Template)

	return drl
}

type DeploymentResponseLite struct {
	ID           int      `json:"id"`
	Name         string   `json:"name"`
	Redeployable bool     `json:"redeployable"`
	ServiceIDs   []string `json:"service_ids"`
	Template     Template `json:"-"`
}

type DeploymentResponseFull struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	Redeployable bool   `json:"redeployable"`
	Status       Status `json:"status"`
}

type Status struct {
	Services Services `json:"services"`
}

type Service struct {
	ID          string `json:"id"`
	ActualState string `json:"actualState"`
}

type Services []Service

type Metadata struct {
	Agent struct {
		Version string `json:"version"`
	} `json:"agent"`
	Adapter interface{} `json:"adapter"`
}
