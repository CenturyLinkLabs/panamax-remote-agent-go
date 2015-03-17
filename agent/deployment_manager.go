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
func (dm DeploymentManager) GetFullDeployment(qid string) (DeploymentResponseFull, error) {
	drl, err := dm.GetDeployment(qid)

	if err != nil {
		return DeploymentResponseFull{}, err
	}

	srvs := make(Services, len(drl.ServiceIDs))
	for i, sID := range drl.ServiceIDs {
		srvc := dm.Adapter.GetService(sID)
		srvs[i] = srvc
	}

	dr := DeploymentResponseFull{
		Name:         drl.Name,
		ID:           drl.ID,
		Redeployable: drl.Redeployable,
		Status:       ServiceStatus{Services: srvs},
	}

	return dr, nil
}

func (dm DeploymentManager) GetDeployment(qid string) (DeploymentResponseLite, error) {
	drl, err := dm.Repo.FindById(qid)
	if err != nil {
		return DeploymentResponseLite{}, err
	}

	return drl, nil
}

func (dm DeploymentManager) DeleteDeployment(dr DeploymentResponseLite) error {
	err := dm.Repo.Remove(dr.ID)

	for _, sID := range dr.ServiceIDs {
		dm.Adapter.DeleteService(sID)
	}

	return err
}

func (dm DeploymentManager) CreateDeployment(body io.Reader) (DeploymentResponseLite, error) {
	deployment := &Deployment{}
	jd := json.NewDecoder(body)
	jd.Decode(deployment)

	ars := dm.Adapter.CreateServices(deployment.MergedImages())

	sIDs := make([]string, len(ars))

	for i, ar := range ars {
		sIDs[i] = ar.ID
	}

	dr, err := dm.Repo.Save(deployment, sIDs)
	if err != nil {
		return DeploymentResponseLite{}, err
	}

	return dr, nil
}
