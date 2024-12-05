package db

import (
	"backend/pkg/images"
	"backend/pkg/quotes"
	"fmt"
	"gorm.io/gorm"
)

type Database interface {
	Connect() error
	GetConnection() (*gorm.DB, error)
	InsertQuote(quote quotes.Quote) error
	FetchQuote(quoteId int) (*quotes.Quote, error)
	FetchAllQuotes() ([]quotes.Quote, error)
	InsertImage(flyer images.Flyer) (uint, error)
	FetchImage(imageId uint) (*images.Flyer, error)
	FetchAllImages() ([]images.Flyer, error)
}

type PostgresDB struct {
	db *gorm.DB
}

func NewPostgresDB() Database {
	return &PostgresDB{}
}

func (p *PostgresDB) Connect() error {
	var err error
	p.db, err = ConnectToDB()
	return err
}

func (p *PostgresDB) GetConnection() (*gorm.DB, error) {
	if p.db == nil {
		return nil, fmt.Errorf("database not connected")
	}
	return p.db, nil
}

func (p *PostgresDB) InsertQuote(quote quotes.Quote) error {
	return p.db.Create(&quote).Error
}

func (p *PostgresDB) FetchQuote(quoteId int) (*quotes.Quote, error) {
	var quote quotes.Quote
	err := p.db.First(&quote, quoteId).Error
	return &quote, err
}

func (p *PostgresDB) FetchAllQuotes() ([]quotes.Quote, error) {
	var quotes []quotes.Quote
	err := p.db.Find(&quotes).Error
	return quotes, err
}

func (p *PostgresDB) InsertImage(flyer images.Flyer) (uint, error) {
	err := p.db.Create(&flyer).Error
	return uint(flyer.Id), err
}

func (p *PostgresDB) FetchImage(imageId uint) (*images.Flyer, error) {
	var flyer images.Flyer
	err := p.db.First(&flyer, imageId).Error
	return &flyer, err
}

func (p *PostgresDB) FetchAllImages() ([]images.Flyer, error) {
	var flyers []images.Flyer
	err := p.db.Find(&flyers).Error
	return flyers, err
} 