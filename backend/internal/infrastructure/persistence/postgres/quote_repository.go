package postgres

import (
	"backend/internal/domain/entity"
	"gorm.io/gorm"
)

type QuoteRepository struct {
	db *gorm.DB
}

func NewQuoteRepository(db *gorm.DB) *QuoteRepository {
	return &QuoteRepository{db: db}
}

func (r *QuoteRepository) Store(quote *entity.Quote) (int, error) {
	result := r.db.Create(quote)
	if result.Error != nil {
		return 0, result.Error
	}
	return quote.Id, nil
}

func (r *QuoteRepository) Update(quote *entity.Quote) error {
	return r.db.Save(quote).Error
}

func (r *QuoteRepository) FindByID(id int) (*entity.Quote, error) {
	var quote entity.Quote
	if err := r.db.First(&quote, id).Error; err != nil {
		return nil, err
	}
	return &quote, nil
}

func (r *QuoteRepository) FindAll() ([]entity.Quote, error) {
	var quotes []entity.Quote
	if err := r.db.Find(&quotes).Error; err != nil {
		return nil, err
	}
	return quotes, nil
} 