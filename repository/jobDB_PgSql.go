package repository

import (
	entity "Atlan_Collect_Ingest/enitity"
	"context"
	"errors"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
)

type JobDbSql struct {
	pool *pgxpool.Pool
}

var TAG_Job = "Sheet Repository"

func NewJobDbSql(pool *pgxpool.Pool) *JobDbSql {
	return &JobDbSql{
		pool: pool,
	}
}

func (r *JobDbSql) Add(pluginCode int) (int, error) {
	log.Printf("%v : Adding Job", TAG_Job)
	var jobId int
	tx, err := r.pool.Begin(context.Background())
	if err != nil {
		return -999, err
	}
	rows, err := tx.Query(context.Background(), "INSERT INTO job_store(job_status,job_status_code,plugin_code) VALUES('In Progress.',202,$1) RETURNING job_id;", pluginCode)
	if err != nil {
		return -999, err
	}
	for rows.Next() {
		err = rows.Scan(&jobId)
		if err != nil {
			return -999, err
		}
	}
	tx.Commit(context.Background())

	return jobId, nil

}

func (r *JobDbSql) Extract(jobId int) (*entity.Job, error) {
	log.Printf("%v : Fetching Job", TAG_Job)
	var status string
	var statusCode int
	var pluginCode int
	tx, err := r.pool.Begin(context.Background())
	if err != nil {
		return nil, err
	}

	rows, err := tx.Query(context.Background(), "SELECT job_status,job_status_code,plugin_code FROM job_store where job_id = $1", jobId)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		err = rows.Scan(&status, &statusCode, &pluginCode)
		if err != nil {
			return nil, err
		}
	}
	tx.Commit(context.Background())

	job := &entity.Job{
		JobId:         jobId,
		JobStatus:     status,
		JobStatusCode: statusCode,
		PluginCode:    pluginCode,
	}
	return job, nil

}

func (r *JobDbSql) Update(jobId, statusCode int, status string) error {
	log.Printf("%s : Updating Job", TAG_Job)
	//UPDATE job_store set job_status = 'Completed.',job_status_code = 200 where job_id = 1
	tx, err := r.pool.Begin(context.Background())
	if err != nil {
		return err
	}
	ct, err := tx.Exec(context.Background(), "UPDATE job_store set job_status = $1,job_status_code = $2 where job_id = $3", status, statusCode, jobId)

	if err != nil {
		log.Printf("err : %v", err.Error())
		return err
	}

	if ct.RowsAffected() < 1 {
		log.Printf("no row affected")
		return errors.New("no row affected")
	}
	tx.Commit(context.Background())
	return nil
}
