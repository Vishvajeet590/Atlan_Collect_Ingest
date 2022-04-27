package repository

import (
	entity "Atlan_Collect_Ingest/enitity"
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
)

type GoogleSheetDbSql struct {
	pool *pgxpool.Pool
}

var TAG_Sheet = "Sheet Repository"

func NewGoogleSheetDbSql(pool *pgxpool.Pool) *GoogleSheetDbSql {
	return &GoogleSheetDbSql{
		pool: pool,
	}
}

func (r *GoogleSheetDbSql) Extract(formId int8) ([]*entity.Responses, error) {
	log.Printf("%s : Fetching Responses for Form %v", TAG_Sheet, formId)

	tx, err := r.pool.Begin(context.Background())
	if err != nil {
		return nil, err
	}

	rows, err := tx.Query(context.Background(), "SELECT question_id,user_id,response from response_store where form_id = $1 ORDER BY user_id ASC ;", formId)
	if err != nil {
		return nil, err
	}
	tx.Commit(context.Background())

	//Creating Slices
	responses := make([]*entity.Responses, 0)
	var prevUser = -999
	var counter = -1
	for rows.Next() {
		var quesId int
		var userId int
		var response string
		err = rows.Scan(&quesId, &userId, &response)
		if err != nil {
			return nil, err
		}
		if userId != prevUser {
			m := make(map[int]string)
			m[quesId] = response
			temp := &entity.Responses{
				UserId:   userId,
				Response: m,
			}
			responses = append(responses, temp)
			prevUser = userId
			counter++
		} else if prevUser == responses[counter].UserId {
			responses[counter].Response[quesId] = response

		}

	}

	rows.Close()

	return responses, nil

}

func (r *GoogleSheetDbSql) QuesIdExtract(formId int8) ([]int, []string, error) {
	//SELECT question_id FROM question_store where form_id = 98
	log.Printf("%s : Fetching Question Ids for Form %v", TAG_Sheet, formId)

	tx, err := r.pool.Begin(context.Background())
	if err != nil {
		return nil, nil, err
	}

	rows, err := tx.Query(context.Background(), "SELECT question_id,question FROM question_store where form_id = $1;", formId)
	if err != nil {
		return nil, nil, err
	}

	quesIds := make([]int, 0)
	questions := make([]string, 0)
	for rows.Next() {
		var qId int
		var ques string
		err = rows.Scan(&qId, &ques)
		if err != nil {
			return nil, nil, err
		}
		quesIds = append(quesIds, qId)
		questions = append(questions, ques)
	}
	rows.Close()
	tx.Commit(context.Background())
	return quesIds, questions, nil
}
