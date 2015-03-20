package agent

import (
	"encoding/json"
	"io"
	"strings"

	"github.com/CenturyLinkLabs/panamax-remote-agent-go/repo"
)

type DeploymentManager struct {
	Repo    repo.DeploymentRepo
	Adapter Adapter
}

func NewDeploymentManager(dRepo repo.DeploymentRepo, ad Adapter) (DeploymentManager, error) {
	return DeploymentManager{
		Repo:    dRepo,
		Adapter: ad,
	}, nil
}

func (dm DeploymentManager) ListDeployments() (DeploymentResponses, error) {
	deps, err := dm.Repo.All()
	if err != nil {
		return DeploymentResponses{}, err
	}

	drs := make(DeploymentResponses, len(deps))

	for i, dep := range deps {
		dr := NewDeploymentResponseLite(
			dep.ID,
			dep.Name,
			dep.Template,
			dep.ServiceIDs,
		)

		drs[i] = *dr
	}

	return drs, nil
}

// TODO: maybe qid should be an int?
func (dm DeploymentManager) GetFullDeployment(qid string) (DeploymentResponseFull, error) {
	dep, err := dm.GetDeployment(qid)

	if err != nil {
		return DeploymentResponseFull{}, err
	}

	srvs := make(Services, len(dep.ServiceIDs))
	for i, sID := range dep.ServiceIDs {
		srvc := dm.Adapter.GetService(sID)
		srvs[i] = srvc
	}

	dr := DeploymentResponseFull{
		Name:         dep.Name,
		ID:           dep.ID,
		Redeployable: dep.Redeployable,
		Status:       ServiceStatus{Services: srvs},
	}

	return dr, nil
}

func (dm DeploymentManager) GetDeployment(qid string) (DeploymentResponseLite, error) {
	dep, err := dm.Repo.FindByID(qid)
	if err != nil {
		return DeploymentResponseLite{}, err
	}

	drl := NewDeploymentResponseLite(
		dep.ID,
		dep.Name,
		dep.Template,
		dep.ServiceIDs,
	)

	return *drl, nil
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

	// decode the template so we can persist it
	b, err := json.Marshal(deployment.Template)
	if err != nil {
		return DeploymentResponseLite{}, err
	}
	template := string(b)

	sb, err := json.Marshal(sIDs)
	sj := string(sb)

	dep := repo.Deployment{
		Name:       deployment.Template.Name,
		Template:   template,
		ServiceIDs: sj,
	}

	if err := dm.Repo.Save(&dep); err != nil {
		return DeploymentResponseLite{}, err
	}

	drl := NewDeploymentResponseLite(
		dep.ID,
		dep.Name,
		dep.Template,
		dep.ServiceIDs,
	)

	return *drl, nil
}

func (dm DeploymentManager) ReDeploy(dr DeploymentResponseLite) (DeploymentResponseLite, error) {

	if err := dm.DeleteDeployment(dr); err != nil {
		return DeploymentResponseLite{}, err
	}

	tBuff, _ := json.Marshal(Deployment{Template: dr.Template})
	tr := strings.NewReader(string(tBuff))

	drl, err := dm.CreateDeployment(tr)
	if err != nil {
		return DeploymentResponseLite{}, err
	}

	return drl, nil
}

func (dm DeploymentManager) FetchMetadata() (Metadata, error) {
	adapterMeta, _ := dm.Adapter.FetchMetadata()

	md := Metadata{
		Agent: struct {
			Version string `json:"version"`
		}{Version: "v1"}, // TODO pull this from a const or ENV or something
		Adapter: adapterMeta,
	}

	return md, nil
}
