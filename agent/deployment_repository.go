package agent

import (
	"database/sql"
	"encoding/json"
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

func (dRepo DeploymentRepo) FindById(qid string) (DeploymentResponseFull, error) {
	var id int
	var name string
	var sids sql.NullString // TODO: make this column NOT NULL
	var template string

	err := dRepo.DB.QueryRow(
		"SELECT id, name, template, service_ids FROM deployments WHERE id = ?",
		qid,
	).Scan(&id, &name, &template, &sids)
	if err != nil {
		// TODO: we could handle ErrNoRows differently, but for now that's just an error to me
		return DeploymentResponseFull{}, err
	}

	dr := DeploymentResponseFull{
		ID:           id,
		Name:         name,
		Redeployable: template != "",
	}

	return dr, nil
}

func (dRepo DeploymentRepo) All() (DeploymentResponses, error) {
	drs := make(DeploymentResponses, 0)

	rows, err := dRepo.DB.Query("SELECT id, name, template, service_ids FROM deployments")
	if err != nil {
		return DeploymentResponses{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var name string
		var sids sql.NullString // TODO: should probably just make this NOT NULL
		var sidc []string
		var template string

		err := rows.Scan(&id, &name, &template, &sids)
		if err != nil {
			return DeploymentResponses{}, err
		}

		if sids.Valid {
			json.Unmarshal([]byte(sids.String), &sidc)
		}

		drs = append(drs, DeploymentResponseLite{
			ID:           id,
			Name:         name,
			ServiceIDs:   sidc,
			Redeployable: template != "",
		})
	}

	if err := rows.Err(); err != nil {
		return DeploymentResponses{}, err
	}

	return drs, err
}

func (dRepo DeploymentRepo) Save(d *Deployment) (DeploymentResponseLite, error) {
	// decode the template so we can persist it
	b, err := json.Marshal(d.Template)
	template := string(b)

	res, err := dRepo.DB.Exec(
		"INSERT INTO deployments (name, template) VALUES (?,?)",
		d.Template.Name,
		template,
	)
	rID, err := res.LastInsertId()
	if err != nil {
		return DeploymentResponseLite{}, err
	}

	dr := DeploymentResponseLite{
		ID:           int(rID),
		Name:         d.Template.Name,
		Redeployable: template != "",
	}

	return dr, nil
}

func (dRepo DeploymentRepo) Remove(qid string) error {
	_, err := dRepo.DB.Exec(
		"DELETE FROM deployments WHERE id = ?",
		qid,
	)
	return err
}
