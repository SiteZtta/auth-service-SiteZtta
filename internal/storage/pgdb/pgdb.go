package pgdb

import (
	"auth-service-SiteZtta/internal/domain/entities"
	"context"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type Storage struct {
	db *sql.DB
}

// New creates a new instance of postgres Storage
func New(connStr string) (*Storage, error) {
	const fn = "auth-service-SiteZtta.internal.storage.pgdb.new"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}
	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}
	storage := &Storage{db: db}
	return storage, nil
}

func (s *Storage) SaveUser(ctx context.Context, user *entities.User) (uid int64, err error) {
	const fn = "auth-service-SiteZtta.internal.storage.pgdb.saveUser"
	err = s.db.QueryRowContext(ctx, "INSERT INTO users (username, email, phone, pass_hash) VALUES ($1, $2, $3, $4) RETURNING id",
		user.Username, user.Email, user.Phone, user.PassHash).Scan(&uid)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", fn, err)
	}
	return uid, nil
}
