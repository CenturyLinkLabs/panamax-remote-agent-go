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

func (dRepo DeploymentRepo) FindByID(qid string) (DeploymentResponseLite, error) {
	var id int
	var name string
	var sids string
	var sidc []string
	var template string

	err := dRepo.DB.QueryRow(
		"SELECT id, name, template, service_ids FROM deployments WHERE id = ?",
		qid,
	).Scan(&id, &name, &template, &sids)
	if err != nil {
		// TODO: we could handle ErrNoRows differently, but for now that's just an error to me
		return DeploymentResponseLite{}, err
	}
	json.Unmarshal([]byte(sids), &sidc)

	tpl := &Template{}
	json.Unmarshal([]byte(template), tpl)

	dr := DeploymentResponseLite{
		ID:           id,
		Name:         name,
		Redeployable: template != "",
		Template:     *tpl,
		ServiceIDs:   sidc,
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
		var sids string
		var sidc []string
		var template string

		err := rows.Scan(&id, &name, &template, &sids)
		if err != nil {
			return DeploymentResponses{}, err
		}

		json.Unmarshal([]byte(sids), &sidc)

		tpl := &Template{}
		json.Unmarshal([]byte(template), tpl)

		drs = append(drs, DeploymentResponseLite{
			ID:           id,
			Name:         name,
			ServiceIDs:   sidc,
			Redeployable: template != "",
			Template:     *tpl,
		})
	}

	if err := rows.Err(); err != nil {
		return DeploymentResponses{}, err
	}

	return drs, err
}

func (dRepo DeploymentRepo) Save(d *Deployment, sIDs []string) (DeploymentResponseLite, error) {
	// decode the template so we can persist it
	b, err := json.Marshal(d.Template)
	template := string(b)

	sb, err := json.Marshal(sIDs)
	sj := string(sb)

	res, err := dRepo.DB.Exec(
		"INSERT INTO deployments (name, template, service_ids) VALUES (?,?,?)",
		d.Template.Name,
		template,
		sj,
	)
	rID, err := res.LastInsertId()
	if err != nil {
		return DeploymentResponseLite{}, err
	}

	dr := DeploymentResponseLite{
		ID:           int(rID),
		Name:         d.Template.Name,
		Redeployable: template != "",
		ServiceIDs:   sIDs,
		Template:     d.Template,
	}

	return dr, nil
}

func (dRepo DeploymentRepo) Remove(id int) error {
	_, err := dRepo.DB.Exec(
		"DELETE FROM deployments WHERE id = ?",
		id,
	)
	return err
}
