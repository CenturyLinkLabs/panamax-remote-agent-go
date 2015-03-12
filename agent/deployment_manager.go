package agent

import (
	"encoding/json"
	"io"
)

type DeploymentManager struct {
	Repo    DeploymentRepo
	Adapter Adapter
}

func NewDeploymentManager(dRepo DeploymentRepo, ad Adapter) (DeploymentManager, error) {
	return DeploymentManager{
		Repo:    dRepo,
		Adapter: ad,
	}, nil
}

func (dm DeploymentManager) ListDeployments() (DeploymentResponses, error) {
	drs, err := dm.Repo.All()
	if err != nil {
		return DeploymentResponses{}, err
	}
	return drs, nil
}

// TODO: maybe qid should be an int?
func (dm DeploymentManager) GetDeployment(qid string) (DeploymentResponseFull, error) {
	dr, err := dm.Repo.FindById(qid)
	if err != nil {
		return DeploymentResponseFull{}, err
	}

	return dr, nil
}

func (dm DeploymentManager) DeleteDeployment(qid string) error {
	err := dm.Repo.Remove(qid)

	return err
}

func (dm DeploymentManager) CreateDeployment(body io.Reader) (DeploymentResponseLite, error) {
	deployment := &Deployment{}
	jd := json.NewDecoder(body)
	jd.Decode(deployment)

	dm.Adapter.CreateServices(deployment.MergedImages())

	dr, err := dm.Repo.Save(deployment)
	if err != nil {
		return DeploymentResponseLite{}, err
	}

	return dr, nil
}
