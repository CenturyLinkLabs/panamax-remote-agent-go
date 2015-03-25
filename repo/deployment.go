package repo

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

func NewPersister(dbPath string) Persister {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		fmt.Println(err)
	}

	return DeploymentRepo{
		DB: db,
	}
}

func (dRepo DeploymentRepo) FindByID(qid int) (Deployment, error) {
	dep := &Deployment{}

	err := dRepo.DB.QueryRow(
		"SELECT id, name, template, service_ids FROM deployments WHERE id = ?",
		qid,
	).Scan(&dep.ID, &dep.Name, &dep.Template, &dep.ServiceIDs)
	if err != nil {
		return Deployment{}, err
	}

	return *dep, nil
}

func (dRepo DeploymentRepo) All() ([]Deployment, error) {
	drs := make([]Deployment, 0)

	rows, err := dRepo.DB.Query("SELECT id, name, template, service_ids FROM deployments")
	if err != nil {
		return []Deployment{}, err
	}
	defer rows.Close()

	for rows.Next() {
		dep := &Deployment{}

		err := rows.Scan(&dep.ID, &dep.Name, &dep.Template, &dep.ServiceIDs)

		if err != nil {
			return []Deployment{}, err
		}

		drs = append(drs, *dep)
	}

	if err := rows.Err(); err != nil {
		return []Deployment{}, err
	}

	return drs, err
}

func (dRepo DeploymentRepo) Save(dep *Deployment) error {
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
