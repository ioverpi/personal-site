package app

import (
	"database/sql"

	"github.com/ioverpi/personal-site/internal/config"
	"github.com/ioverpi/personal-site/internal/database"
)

type App struct {
	DB     *sql.DB
	Config *config.Config
}

func New(cfg *config.Config) (*App, error) {
	db, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		return nil, err
	}

	// Run pending migrations
	if err := database.Migrate(db); err != nil {
		db.Close()
		return nil, err
	}

	return &App{
		DB:     db,
		Config: cfg,
	}, nil
}

func (a *App) Close() {
	if a.DB != nil {
		a.DB.Close()
	}
}
