package agent

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

const (
	dbDriver = "sqlite"
)

type DeploymentRepo struct {
	DB *sql.DB
}

func NewDeploymentRepo(dbPath string) DeploymentRepo {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		fmt.Println(err)
	}
	//TODO: do I ever need to call close?

	return DeploymentRepo{
		DB: db,
	}
}

func (dRepo DeploymentRepo) FindByID(qid string) (Rdeployment, error) {
	dep := &Rdeployment{}

	err := dRepo.DB.QueryRow(
		"SELECT id, name, template, service_ids FROM deployments WHERE id = ?",
		qid,
	).Scan(&dep.ID, &dep.Name, &dep.Template, &dep.ServiceIDs)
	if err != nil {
		// TODO: we could handle ErrNoRows differently, but for now that's just an error to me
		return Rdeployment{}, err
	}

	return *dep, nil
}

func (dRepo DeploymentRepo) All() ([]Rdeployment, error) {
	drs := make([]Rdeployment, 0)

	rows, err := dRepo.DB.Query("SELECT id, name, template, service_ids FROM deployments")
	if err != nil {
		return []Rdeployment{}, err
	}
	defer rows.Close()

	for rows.Next() {
		dep := &Rdeployment{}

		err := rows.Scan(&dep.ID, &dep.Name, &dep.Template, &dep.ServiceIDs)

		if err != nil {
			return []Rdeployment{}, err
		}

		drs = append(drs, *dep)
	}

	if err := rows.Err(); err != nil {
		return []Rdeployment{}, err
	}

	return drs, err
}

func (dRepo DeploymentRepo) Save(dep *Rdeployment) error {
	res, err := dRepo.DB.Exec(
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

func (dRepo DeploymentRepo) Remove(id int) error {
	_, err := dRepo.DB.Exec(
		"DELETE FROM deployments WHERE id = ?",
		id,
	)
	return err
}

// eventually this can be Repo.Deployment or something
type Rdeployment struct {
	ID         int
	Name       string
	ServiceIDs string
	Template   string
}
