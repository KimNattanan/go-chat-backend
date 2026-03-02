// Package postgres implements postgres connection.
package postgres

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
	_defaultConnAttempts = 10
	_defaultConnTimeout  = time.Second
)

// Postgres -.
type Postgres struct {
	connAttempts int
	connTimeout  time.Duration
	db           *gorm.DB
	sqlDB        *sql.DB
}

// New -.
func New(dsn string, opts ...Option) (*Postgres, error) {
	pg := &Postgres{
		connAttempts: _defaultConnAttempts,
		connTimeout:  _defaultConnTimeout,
	}

	// Custom options
	for _, opt := range opts {
		opt(pg)
	}

	var err error

	for pg.connAttempts > 0 {
		pg.db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err == nil {
			break
		}

		log.Printf("Postgres is trying to connect, attempts left: %d", pg.connAttempts)

		time.Sleep(pg.connTimeout)

		pg.connAttempts--
	}

	if err != nil {
		return nil, fmt.Errorf("postgres - NewPostgres - connAttempts == 0: %w", err)
	}

	pg.sqlDB, err = pg.db.DB()
	if err != nil {
		return nil, fmt.Errorf("postgres - NewPostgres - pg.db.DB: %w", err)
	}
	return pg, nil
}

// Close -.
func (p *Postgres) Close() {
	if p.sqlDB != nil {
		p.sqlDB.Close()
	}
}
