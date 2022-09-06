package demo

type Service struct {
	repository *Repository
}

func NewService() *Service {
	repository := NewRepository()
	return &Service{
		repository: repository,
	}
}

func (service *Service) GetUsers() []UserModel {
	ids := service.repository.GetUserIds()
	return service.repository.GetUserByIds(ids)
}
