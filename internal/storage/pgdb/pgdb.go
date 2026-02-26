package pgdb

import (
	"auth-service-SiteZtta/internal/domain/entities"
	"auth-service-SiteZtta/internal/storage"
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/lib/pq"
	_ "github.com/lib/pq"
)

type Storage struct {
	db *sql.DB
}

const (
	usersTable = "users"
)

// New creates a new instance of postgres Storage
func New(connStr string) (*Storage, error) {
	const fn = "auth-service-SiteZtta.internal.storage.pgdb.new"
	fmt.Printf("connStr %s\n", connStr)
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

// SaveUser saves a new user to the database
func (s *Storage) SaveUser(ctx context.Context, user *entities.User) (uid int64, err error) {
	const fn = "auth-service-SiteZtta.internal.storage.pgdb.saveUser"
	stmt, err := s.db.PrepareContext(ctx, fmt.Sprintf("INSERT INTO %s (username, email, phone, pass_hash) VALUES ($1, $2, $3, $4) RETURNING id", usersTable))
	if err != nil {
		return 0, fmt.Errorf("%s: %w", fn, err)
	}
	rows, err := stmt.QueryContext(ctx, user.Username, user.Email, user.Phone, user.PassHash)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" { // unique_violation
			return 0, fmt.Errorf("%s: %w", fn, storage.ErrUserExists)
		}
		return 0, fmt.Errorf("%s: %w", fn, err)
	}
	defer rows.Close()
	if rows.Next() {
		if err = rows.Scan(&uid); err != nil {
			return 0, fmt.Errorf("%s: %w", fn, err)
		}
	}
	return uid, nil
}

// GerUserByEmail gets a user by email
func (s *Storage) GetUserByEmail(ctx context.Context, email string) (user *entities.User, err error) {
	const fn = "auth-service-SiteZtta.internal.storage.pgdb.getUserByEmail"
	stmt, err := s.db.PrepareContext(ctx, fmt.Sprintf("SELECT id, username, email, phone, pass_hash, role FROM %s WHERE email = $1", usersTable))
	if err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}
	row := stmt.QueryRowContext(ctx, email)
	user = &entities.User{}
	err = row.Scan(&user.ID, &user.Username, &user.Email, &user.Phone, &user.PassHash, &user.Role)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", fn, storage.ErrUserNotFound)
		}
		return nil, fmt.Errorf("%s: %w", fn, err)
	}
	return user, nil
}

// GetUserByUsername gets a user by username
func (s *Storage) GetUserByUsername(ctx context.Context, username string) (user *entities.User, err error) {
	const fn = "auth-service-SiteZtta.internal.storage.pgdb.getUserByUsername"
	stmt, err := s.db.PrepareContext(ctx, fmt.Sprintf("SELECT id, username, email, phone, pass_hash, role FROM %s WHERE username = $1", usersTable))
	if err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}
	rows := stmt.QueryRowContext(ctx, username)
	user = &entities.User{}
	err = rows.Scan(&user.ID, &user.Username, &user.Email, &user.Phone, &user.PassHash, &user.Role)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", fn, storage.ErrUserNotFound)
		}
		return nil, fmt.Errorf("%s: %w", fn, err)
	}
	return user, nil
}
