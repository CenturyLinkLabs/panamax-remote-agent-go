package api

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/CenturyLinkLabs/panamax-remote-agent-go/agent"
	"github.com/gorilla/mux"
)

func listDeployments(dm agent.DeploymentManager, w http.ResponseWriter, r *http.Request) {
	drs, err := dm.ListDeployments()
	if err != nil {
		log.Fatal(err)
	}

	if err := json.NewEncoder(w).Encode(drs); err != nil {
		log.Fatal(err)
	}
}

func createDeployment(dm agent.DeploymentManager, w http.ResponseWriter, r *http.Request) {
	depB := &agent.DeploymentBlueprint{}
	jd := json.NewDecoder(r.Body)
	if err := jd.Decode(depB); err != nil {
		log.Fatal(err)
	}

	dr, err := dm.CreateDeployment(*depB)
	if err != nil {
		log.Fatal(err)
	}

	drj, errr := json.Marshal(dr)
	if errr != nil {
		log.Fatal(errr)
	}

	w.WriteHeader(http.StatusCreated)
	_, err = w.Write(drj)
	if err != nil {
		log.Fatal(err)
	}
}

func deleteDeployment(dm agent.DeploymentManager, w http.ResponseWriter, r *http.Request) {
	if err := dm.DeleteDeployment(idFromQuery(mux.Vars(r))); err != nil {
		log.Fatal(err)
	}

	w.WriteHeader(http.StatusNoContent)
}

func showDeployment(dm agent.DeploymentManager, w http.ResponseWriter, r *http.Request) {
	dr, err := dm.GetFullDeployment(idFromQuery(mux.Vars(r)))
	if err != nil {
		log.Fatal(err)
	}

	if err := json.NewEncoder(w).Encode(dr); err != nil {
		log.Fatal(err)
	}
}

func reDeploy(dm agent.DeploymentManager, w http.ResponseWriter, r *http.Request) {
	dr, err := dm.ReDeploy(idFromQuery(mux.Vars(r)))
	if err != nil {
		log.Fatal(err)
	}

	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(dr); err != nil {
		log.Fatal(err)
	}
}

func metadata(dm agent.DeploymentManager, w http.ResponseWriter, r *http.Request) {
	md, _ := dm.FetchMetadata()

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(md); err != nil {
		log.Fatal(err)
	}
}

func idFromQuery(uvars map[string]string) int {
	qID, err := strconv.Atoi(uvars["id"])

	if err != nil {
		log.Fatal(err)
	}

	return qID
}
