package agent

import (
	"encoding/json"
)

// A DeploymentBlueprint is the top level entity, containing all the
// necessary bits to kick off a deployment.
type DeploymentBlueprint struct {
	Override Template `json:"override,omitempty"`
	Template Template `json:"template,omitempty"`
}

// MergedImages merges the Override on top of the Template, returning the
// resulting merged Images to be used for deployment.
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

// A Template is the 2nd level entity in the DeploymentBlueprint scheme.
// It contains all the necessary information for a deployment post override logic.
type Template struct {
	Name   string  `json:"name,omitempty"`
	Images []Image `json:"images,omitempty"`
}

// An Image ultimately represents the deployed Docker image.
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

// MarshalJSON is used to strip out empty/default value structs when
// marshalling images to JSON.
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

// Environment represents each environment variable that will be passed
// to the Docker run command.
type Environment struct {
	Variable string `json:"variable,omitempty"`
	Value    string `json:"value,omitempty"`
}

// Link represents each Link that will be passed
// to the Docker run command.
type Link struct {
	Service string `json:"service,omitempty"`
	Alias   string `json:"alias,omitempty"`
}

// Port represents each Port mapping that will be passed
// to the Docker run command.
type Port struct {
	HostPort      int `json:"hostPort,omitempty"`
	ContainerPort int `json:"containerPort,omitempty"`
}

// DeploymentSettings contains orchestrator specific information
// to be used when deploying an application.
type DeploymentSettings struct {
	Count int `json:"count,omitempty"`
}

// A Volume represents each Volume mapping that will be passed
// to the Docker run command.
type Volume struct {
	ContainerPath string `json:"containerPath"`
	HostPath      string `json:"hostPath"`
}

// DeploymentResponseLite is the minimal representation of a Deployment
// typically used for listings, etc.
type DeploymentResponseLite struct {
	ID           int      `json:"id"`
	Name         string   `json:"name"`
	Redeployable bool     `json:"redeployable"`
	ServiceIDs   []string `json:"service_ids"`
}

// DeploymentResponseFull is the robust representation of a Deployment
// typically used for the return value of an individual deployment, etc.
type DeploymentResponseFull struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	Redeployable bool   `json:"redeployable"`
	Status       Status `json:"status"`
}

// Status contains information for health of each service in a Deployment.
type Status struct {
	Services []Service `json:"services"`
}

// Service represents each service in a Deployment and contains the ID,
// as well as its state.
type Service struct {
	ID          string `json:"id"`
	ActualState string `json:"actualState"`
}

// Metadata contains general meta data for both the Agent and the Adapter.
type Metadata struct {
	Agent struct {
		Version string `json:"version"`
	} `json:"agent"`
	Adapter interface{} `json:"adapter"`
}
