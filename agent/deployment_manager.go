package agent

import (
	"bytes"
	"encoding/json"

	"github.com/CenturyLinkLabs/panamax-remote-agent-go/adapter"
	"github.com/CenturyLinkLabs/panamax-remote-agent-go/repo"
)

type DeploymentManager struct {
	Repo   repo.Persister
	Client adapter.Client
}

func MakeDeploymentManager(p repo.Persister, c adapter.Client) DeploymentManager {
	return DeploymentManager{
		Repo:   p,
		Client: c,
	}
}

func (dm DeploymentManager) ListDeployments() (DeploymentResponses, error) {
	deps, err := dm.Repo.All()
	if err != nil {
		return DeploymentResponses{}, err
	}

	drs := make(DeploymentResponses, len(deps))

	for i, dep := range deps {
		dr := deploymentResponseLiteFromRawValues(
			dep.ID,
			dep.Name,
			dep.Template,
			dep.ServiceIDs,
		)

		drs[i] = dr
	}

	return drs, nil
}

func (dm DeploymentManager) GetFullDeployment(qid int) (DeploymentResponseFull, error) {
	dep, err := dm.GetDeployment(qid)

	if err != nil {
		return DeploymentResponseFull{}, err
	}

	//TODO: maybe use a constructor for all this.
	as := make(Services, len(dep.ServiceIDs))
	for i, sID := range dep.ServiceIDs {
		srvc := dm.Client.GetService(sID)
		as[i] = Service{
			ID:          srvc.ID,
			ActualState: srvc.ActualState,
		}
	}

	dr := DeploymentResponseFull{
		Name:         dep.Name,
		ID:           dep.ID,
		Redeployable: dep.Redeployable,
		Status:       Status{Services: as},
	}

	return dr, nil
}

func (dm DeploymentManager) GetDeployment(qid int) (DeploymentResponseLite, error) {
	dep, err := dm.Repo.FindByID(qid)
	if err != nil {
		return DeploymentResponseLite{}, err
	}

	drl := deploymentResponseLiteFromRawValues(
		dep.ID,
		dep.Name,
		dep.Template,
		dep.ServiceIDs,
	)

	return drl, nil
}

func (dm DeploymentManager) DeleteDeployment(qID int) error {
	dep, err := dm.Repo.FindByID(qID)

	if err != nil {
		return err
	}

	var sIDs []string
	json.Unmarshal([]byte(dep.ServiceIDs), &sIDs)

	for _, sID := range sIDs {
		dm.Client.DeleteService(sID)
	}

	dm.Repo.Remove(qID)

	return err
}

func (dm DeploymentManager) CreateDeployment(depB DeploymentBlueprint) (DeploymentResponseLite, error) {

	mImgs := depB.MergedImages()

	buf := &bytes.Buffer{}
	if err := json.NewEncoder(buf).Encode(mImgs); err != nil {
		return DeploymentResponseLite{}, err
	}

	as := dm.Client.CreateServices(buf)

	tn := depB.Template.Name
	dep, err := makeRepoDeployment(tn, mImgs, as)
	if err != nil {
		return DeploymentResponseLite{}, err
	}

	if err := dm.Repo.Save(&dep); err != nil {
		return DeploymentResponseLite{}, err
	}

	drl := deploymentResponseLiteFromRawValues(
		dep.ID,
		dep.Name,
		dep.Template,
		dep.ServiceIDs,
	)

	return drl, nil
}

func (dm DeploymentManager) ReDeploy(ID int) (DeploymentResponseLite, error) {

	dep, err := dm.Repo.FindByID(ID)

	dr := deploymentResponseLiteFromRawValues(
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
	adapterMeta, _ := dm.Client.FetchMetadata()

	md := Metadata{
		Agent: struct {
			Version string `json:"version"`
		}{Version: "v1"}, // TODO pull this from a const or ENV or something
		Adapter: adapterMeta,
	}

	return md, nil
}

func makeRepoDeployment(tn string, mImgs []Image, as []adapter.Service) (repo.Deployment, error) {
	ts, err := stringifyTemplate(tn, mImgs)
	ss, err := stringifyServiceIDs(as)

	if err != nil {
		return repo.Deployment{}, err
	}

	return repo.Deployment{
		Name:       tn,
		Template:   ts,
		ServiceIDs: ss,
	}, nil
}

func stringifyTemplate(tn string, imgs []Image) (string, error) {
	mt := Template{
		Name:   tn,
		Images: imgs,
	}
	b, err := json.Marshal(mt)

	return string(b), err
}

func stringifyServiceIDs(as []adapter.Service) (string, error) {
	sIDs := make([]string, len(as))

	for i, ar := range as {
		sIDs[i] = ar.ID
	}

	sb, err := json.Marshal(sIDs)

	return string(sb), err
}

func deploymentResponseLiteFromRawValues(id int, nm string, tpl string, sids string) DeploymentResponseLite {
	drl := &DeploymentResponseLite{
		ID:           id,
		Name:         nm,
		Redeployable: tpl != "",
	}
	json.Unmarshal([]byte(sids), &drl.ServiceIDs)
	json.Unmarshal([]byte(tpl), &drl.Template)

	return *drl
}
