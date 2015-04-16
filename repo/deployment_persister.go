package repo

import (
	"database/sql"

	"github.com/CenturyLinkLabs/panamax-remote-agent-go/agent"

	// The sql.Open command references the driver name
	_ "github.com/mattn/go-sqlite3"
)

const (
	dbDriver = "sqlite"
)

// MakeDeploymentStore returns a new deploymentPersister as an agent.DeploymentStore.
// agent.DeploymentStore interface.
func MakeDeploymentStore(dbPath string) (agent.DeploymentStore, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return deploymentPersister{}, err
	}

	return deploymentPersister{
		db: db,
	}, nil
}

type deploymentPersister struct {
	db *sql.DB
}

func (p deploymentPersister) FindByID(qid int) (agent.Deployment, error) {
	dep := &agent.Deployment{}

	err := p.db.QueryRow(
		"SELECT id, name, template, service_ids FROM deployments WHERE id = ?",
		qid,
	).Scan(&dep.ID, &dep.Name, &dep.Template, &dep.ServiceIDs)
	if err != nil {
		return agent.Deployment{}, err
	}

	return *dep, nil
}

func (p deploymentPersister) All() ([]agent.Deployment, error) {
	var ds []agent.Deployment

	rows, err := p.db.Query("SELECT id, name, template, service_ids FROM deployments")
	if err != nil {
		return []agent.Deployment{}, err
	}
	defer rows.Close()

	for rows.Next() {
		dep := &agent.Deployment{}

		err := rows.Scan(&dep.ID, &dep.Name, &dep.Template, &dep.ServiceIDs)

		if err != nil {
			return []agent.Deployment{}, err
		}

		ds = append(ds, *dep)
	}

	if err := rows.Err(); err != nil {
		return []agent.Deployment{}, err
	}

	return ds, err
}

func (p deploymentPersister) Save(dep *agent.Deployment) error {
	res, err := p.db.Exec(
		"INSERT INTO deployments (name, template, service_ids) VALUES (?,?,?)",
		dep.Name,
		dep.Template,
		dep.ServiceIDs,
	)
	rID, err := res.LastInsertId()
	if err != nil {
		return err
	}

	dep.ID = int(rID)

	return nil
}

func (p deploymentPersister) Remove(id int) error {
	_, err := p.db.Exec(
		"DELETE FROM deployments WHERE id = ?",
		id,
	)
	return err
}
