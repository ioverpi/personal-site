package services

import (
	"database/sql"
	"math/rand"

	"github.com/ioverpi/personal-site/internal/app"
	"github.com/ioverpi/personal-site/internal/models"
)

type QuotesService struct {
	app *app.App
}

func NewQuotesService(app *app.App) *QuotesService {
	return &QuotesService{app: app}
}

func (s *QuotesService) GetAllQuotes() ([]models.Quote, error) {
	rows, err := s.app.DB.Query(`
		SELECT id, content, author, is_own, created_at
		FROM quotes
		ORDER BY created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanQuotes(rows)
}

func (s *QuotesService) GetQuoteByID(id int) (*models.Quote, error) {
	var quote models.Quote
	err := s.app.DB.QueryRow(`
		SELECT id, content, author, is_own, created_at
		FROM quotes
		WHERE id = $1
	`, id).Scan(
		&quote.ID, &quote.Content, &quote.Author, &quote.IsOwn, &quote.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &quote, nil
}

func (s *QuotesService) GetRandomQuote() (*models.Quote, error) {
	var quote models.Quote
	err := s.app.DB.QueryRow(`
		SELECT id, content, author, is_own, created_at
		FROM quotes
		ORDER BY RANDOM()
		LIMIT 1
	`).Scan(
		&quote.ID, &quote.Content, &quote.Author, &quote.IsOwn, &quote.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &quote, nil
}

func scanQuotes(rows *sql.Rows) ([]models.Quote, error) {
	var quotes []models.Quote
	for rows.Next() {
		var quote models.Quote
		err := rows.Scan(
			&quote.ID, &quote.Content, &quote.Author, &quote.IsOwn, &quote.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		quotes = append(quotes, quote)
	}
	return quotes, rows.Err()
}

// Ensure rand is available for potential future use
var _ = rand.Int
