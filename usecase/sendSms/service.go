package sendSms

type Service struct {
	repo Repository
}

func NewService(r Repository) *Service {
	return &Service{
		repo: r,
	}
}

func (s *Service) SendSMS() error {

	for i := 0; i < 10; i++ {
		println("Sending SMS is Same as reading resonse column and integrating twilio or similar api")

	}
	return nil
}
