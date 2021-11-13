package coinprice

import (
	"test-server/models"

	"gorm.io/gorm"
)

// Repository Interface
type Repository interface {
	FindLast() (models.CoinPrice, error)
	Save(models.CoinPrice) (models.CoinPrice, error)
	FindByTime(time string) (models.CoinPrice, error)
	FindBetween(from string, to string) ([]models.CoinPrice, error)
	BatchInsert(input []models.CoinPrice) error
	Truncate() error
}

//repository struct
type repository struct {
	db *gorm.DB
}

// NewRepo is the single instance repo that is being created.
func NewRepo(db *gorm.DB) Repository {
	return &repository{
		db: db,
	}
}

// Find - Find Last price
func (r *repository) FindLast() (models.CoinPrice, error) {
	var data models.CoinPrice
	result := r.db.Last(&data)
	if result.Error != nil {
		return data, result.Error
	}
	return data, nil
}

// Save - new rate
func (r *repository) Save(input models.CoinPrice) (models.CoinPrice, error) {
	var data = input
	result := r.db.Save(&data)
	if result.Error != nil {
		return data, result.Error
	}
	return data, nil
}

// Batch - insert Rates
func (r *repository) BatchInsert(input []models.CoinPrice) error {

	result := r.db.CreateInBatches(input, 100)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// Find rate by timestamp
func (r *repository) FindByTime(time string) (models.CoinPrice, error) {
	var data models.CoinPrice

	result := r.db.Where("Time = ?", time).Find(&data)

	if result.Error != nil {
		return data, result.Error
	}
	return data, nil
}

// Find rates between given timestamp
func (r *repository) FindBetween(from string, to string) ([]models.CoinPrice, error) {
	var data []models.CoinPrice

	result := r.db.Where("Time BETWEEN ? AND ?", from, to).Find(&data)

	if result.Error != nil {
		return nil, result.Error
	}
	return data, nil
}

// Delete all data from coin_prices table
func (r *repository) Truncate() error {

	result := r.db.Exec("DELETE FROM coin_prices")
	if result.Error != nil {
		return result.Error
	}
	return nil
}
