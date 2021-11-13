package coinprice

import (
	"test-server/models"
)

type Service interface {
	FindLast() (models.CoinPrice, error)
	Save(models.CoinPrice) (models.CoinPrice, error)
	FindByTime(time string) (models.CoinPrice, error)
	FindBetween(from string, to string) ([]models.CoinPrice, error)
	BatchInsert(input []models.CoinPrice) error
	Truncate() error
}

// service struct
type service struct {
	repository Repository
}

// NewService is the single instance service that is being created.
func NewService(r Repository) Service {
	return &service{
		repository: r,
	}
}

func (s *service) FindLast() (models.CoinPrice, error) {
	return s.repository.FindLast()
}

func (s *service) Save(input models.CoinPrice) (models.CoinPrice, error) {
	return s.repository.Save(input)
}

func (s *service) FindByTime(time string) (models.CoinPrice, error) {
	return s.repository.FindByTime(time)
}

func (s *service) FindBetween(from string, to string) ([]models.CoinPrice, error) {
	return s.repository.FindBetween(from, to)
}

func (s *service) BatchInsert(input []models.CoinPrice) error {
	return s.repository.BatchInsert(input)
}

func (s *service) Truncate() error {
	return s.repository.Truncate()
}
