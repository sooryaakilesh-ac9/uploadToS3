package repository

import "backend/internal/domain/entity"

type QuoteRepository interface {
    Store(quote *entity.Quote) (int, error)
    Update(quote *entity.Quote) error
    FindByID(id int) (*entity.Quote, error)
    FindAll() ([]entity.Quote, error)
} 