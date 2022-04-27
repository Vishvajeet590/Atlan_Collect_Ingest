package sendSms

import entity "Atlan_Collect_Ingest/enitity"

type Reader interface {
	Extracts(formId int8) ([]*entity.Responses, error)
	QuesIdExtracts(formId int8) ([]int, []string, error)
}

type Repository interface {
	Reader
}

type Usecase interface {
	SendSMS() error
}
