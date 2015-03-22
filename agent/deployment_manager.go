package agent

import (
	"bytes"
	"encoding/json"

	"github.com/CenturyLinkLabs/panamax-remote-agent-go/adapter"
	"github.com/CenturyLinkLabs/panamax-remote-agent-go/repo"
)

type DeploymentManager struct {
	Repo          repo.DeploymentRepo //TODO, this should probably be an interface
	AdapterClient adapter.Client
}

func NewDeploymentManager(dRepo repo.DeploymentRepo, ad adapter.Client) (DeploymentManager, error) {
	return DeploymentManager{
		Repo:          dRepo,
		AdapterClient: ad,
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

func (dm DeploymentManager) GetFullDeployment(qid int) (DeploymentResponseFull, error) {
	dep, err := dm.GetDeployment(qid)

	if err != nil {
		return DeploymentResponseFull{}, err
	}

	srvs := make(Services, len(dep.ServiceIDs))
	for i, sID := range dep.ServiceIDs {
		srvc := dm.AdapterClient.GetService(sID)
		srvs[i] = Service{
			ID:          srvc.ID,
			ActualState: srvc.ActualState,
		}
	}

	dr := DeploymentResponseFull{
		Name:         dep.Name,
		ID:           dep.ID,
		Redeployable: dep.Redeployable,
		Status:       ServiceStatus{Services: srvs},
	}

	return dr, nil
}

func (dm DeploymentManager) GetDeployment(qid int) (DeploymentResponseLite, error) {
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

func (dm DeploymentManager) DeleteDeployment(qID int) error {
	dep, err := dm.Repo.FindByID(qID)

	if err != nil {
		return err
	}

	var sIDs []string
	json.Unmarshal([]byte(dep.ServiceIDs), &sIDs)

	for _, sID := range sIDs {
		dm.AdapterClient.DeleteService(sID)
	}

	if err := dm.Repo.Remove(qID); err != nil {
		return err
	}

	return err
}

func (dm DeploymentManager) CreateDeployment(depB DeploymentBlueprint) (DeploymentResponseLite, error) {

	imgs := depB.MergedImages()

	buf := new(bytes.Buffer)
	if err := json.NewEncoder(buf).Encode(imgs); err != nil {
		panic(err)
	}

	ars := dm.AdapterClient.CreateServices(buf)

	sIDs := make([]string, len(ars))

	for i, ar := range ars {
		sIDs[i] = ar.ID
	}

	// decode the template so we can persist it
	b, err := json.Marshal(depB.Template)
	if err != nil {
		return DeploymentResponseLite{}, err
	}
	template := string(b)

	sb, err := json.Marshal(sIDs)
	sj := string(sb)

	dep := repo.Deployment{
		Name:       depB.Template.Name,
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

func (dm DeploymentManager) ReDeploy(ID int) (DeploymentResponseLite, error) {

	dep, err := dm.Repo.FindByID(ID)

	dr := NewDeploymentResponseLite(
		dep.ID,
		dep.Name,
		dep.Template,
		dep.ServiceIDs,
	)

	if err := dm.DeleteDeployment(ID); err != nil {
		return DeploymentResponseLite{}, err
	}

	drl, err := dm.CreateDeployment(DeploymentBlueprint{Template: dr.Template})
	if err != nil {
		return DeploymentResponseLite{}, err
	}

	return drl, nil
}

func (dm DeploymentManager) FetchMetadata() (Metadata, error) {
	adapterMeta, _ := dm.AdapterClient.FetchMetadata()

	md := Metadata{
		Agent: struct {
			Version string `json:"version"`
		}{Version: "v1"}, // TODO pull this from a const or ENV or something
		Adapter: adapterMeta,
	}

	return md, nil
}
