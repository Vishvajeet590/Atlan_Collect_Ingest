package jobStatus

import entity "Atlan_Collect_Ingest/enitity"

type Service struct {
	repo Repository
}

func NewService(r Repository) *Service {
	return &Service{
		repo: r,
	}
}

func (s *Service) AddJob(pluginCode int) (int, error) {

	jobId, err := s.repo.Add(pluginCode)
	if err != nil {
		return -999, err
	}
	return jobId, nil

}
func (s *Service) UpdateJob(jobId, statusCode int, status string) error {

	err := s.repo.Update(jobId, statusCode, status)
	if err != nil {
		return err
	}
	return nil

}
func (s *Service) GetJob(jobId int) (*entity.Job, error) {
	job, err := s.repo.Extract(jobId)
	if err != nil {
		return nil, err
	}
	return job, nil
}
