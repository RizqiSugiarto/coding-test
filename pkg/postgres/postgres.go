package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/Masterminds/squirrel"
	// Import postgres driver.
	_ "github.com/lib/pq"
)

const (
	_defaultMaxPoolSize     = 1
	_defaultConnAttempts    = 10
	_defaultConnTimeout     = time.Second
	_defaultPingTimeout     = 2 * time.Second
	_defaultConnMaxLifetime = 5 * time.Minute
)

// Postgres represents a PostgreSQL connection handler.
type Postgres struct {
	maxPoolSize  int
	connAttempts int
	connTimeout  time.Duration

	Builder squirrel.StatementBuilderType
	DB      *sql.DB
}

// New initializes a new Postgres connection.
func New(url string, opts ...Option) (*Postgres, error) {
	pg := &Postgres{
		maxPoolSize:  _defaultMaxPoolSize,
		connAttempts: _defaultConnAttempts,
		connTimeout:  _defaultConnTimeout,
	}

	for _, opt := range opts {
		opt(pg)
	}

	pg.Builder = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	var db *sql.DB

	var err error

	for attempt := 1; attempt <= pg.connAttempts; attempt++ {
		db, err = sql.Open("postgres", url)
		if err != nil {
			log.Printf("Postgres connection failed (attempt %d/%d): %v", attempt, pg.connAttempts, err)
			time.Sleep(pg.connTimeout)

			continue
		}

		ctx, cancel := context.WithTimeout(context.Background(), _defaultPingTimeout)
		err = db.PingContext(ctx)

		cancel()

		if err != nil {
			log.Printf("Postgres ping failed (attempt %d/%d): %v", attempt, pg.connAttempts, err)

			_ = db.Close()

			time.Sleep(pg.connTimeout)

			continue
		}

		// Success — stop retrying
		break
	}

	if err != nil {
		return nil, fmt.Errorf("postgres - NewPostgres - all connection attempts failed: %w", err)
	}

	db.SetMaxOpenConns(pg.maxPoolSize)
	db.SetConnMaxLifetime(_defaultConnMaxLifetime)

	pg.DB = db

	return pg, nil
}

// Close closes the database connection.
func (p *Postgres) Close() {
	if p.DB != nil {
		_ = p.DB.Close()
	}
}
