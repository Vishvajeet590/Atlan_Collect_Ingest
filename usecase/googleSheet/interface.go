package googleSheet

import entity "Atlan_Collect_Ingest/enitity"

type Reader interface {
	Extract(formId int8) ([]*entity.Responses, error)
	QuesIdExtract(formId int8) ([]int, []string, error)
}

type Repository interface {
	Reader
}

type Usecase interface {
	AddToSheet(formId int8, oAuthCode string) ([]*entity.Responses, error)
}
