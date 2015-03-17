package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/CenturyLinkLabs/panamax-remote-agent-go/agent"
	"github.com/gorilla/mux"
)

func ListDeployments(dm agent.DeploymentManager, w http.ResponseWriter, r *http.Request) {
	drs, err := dm.ListDeployments()
	if err != nil {
		log.Fatal(err)
	}

	json.NewEncoder(w).Encode(drs)
}

func CreateDeployment(dm agent.DeploymentManager, w http.ResponseWriter, r *http.Request) {
	dr, err := dm.CreateDeployment(r.Body)
	if err != nil {
		log.Fatal(err)
	}

	drj, errr := json.Marshal(dr)
	if errr != nil {
		log.Fatal(errr)
	}

	w.WriteHeader(http.StatusCreated)
	w.Write(drj)
}

func DeleteDeployment(dm agent.DeploymentManager, w http.ResponseWriter, r *http.Request) {
	uvars := mux.Vars(r)
	dr, err := dm.GetDeployment(uvars["id"])
	if err != nil {
		log.Fatal(err)
	}

	if err := dm.DeleteDeployment(dr); err != nil {
		log.Fatal(err)
	}

	w.WriteHeader(http.StatusNoContent)
}

func ShowDeployment(dm agent.DeploymentManager, w http.ResponseWriter, r *http.Request) {
	uvars := mux.Vars(r)
	dr, err := dm.GetFullDeployment(uvars["id"])
	if err != nil {
		log.Fatal(err)
	}

	json.NewEncoder(w).Encode(dr)
}

func ReDeploy(dm agent.DeploymentManager, w http.ResponseWriter, r *http.Request) {
	uvars := mux.Vars(r)
	dr, err := dm.GetDeployment(uvars["id"])
	if err != nil {
		log.Fatal(err)
	}

	rdr, err := dm.ReDeploy(dr)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(rdr)
}
