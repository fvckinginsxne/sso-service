package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/lib/pq"
	_ "github.com/lib/pq"

	"auth/internal/domain/models"
	"auth/internal/storage"
)

type Storage struct {
	db *sql.DB
}

func New(dsn string) (*Storage, error) {
	const op = "storage.postgres.New"

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) SaveUser(ctx context.Context,
	email string,
	passHash []byte,
) error {
	const op = "storage.postgres.SaveUser"

	stmt, err := s.db.PrepareContext(ctx, "INSERT INTO users (email, pass_hash) VALUES ($1, $2)")
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	defer stmt.Close()

	if _, err := stmt.ExecContext(ctx, email, passHash); err != nil {
		var pgErr *pq.Error
		if errors.As(err, &pgErr) && pgErr.Code == storage.UniqueViolationCode {
			return fmt.Errorf("%s: %w", op, storage.ErrUserExists)
		}

		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) User(ctx context.Context, email string) (*models.User, error) {
	const op = "storage.postgres.User"

	stmt, err := s.db.PrepareContext(ctx, "SELECT id, email, pass_hash FROM users WHERE email = $1")
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer stmt.Close()

	var user models.User
	if err := stmt.QueryRowContext(ctx, email).Scan(&user.ID, &user.Email, &user.PassHash); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", op, storage.ErrUserNotFound)
		}

		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &user, nil
}

func (s *Storage) Close() error {
	const op = "storage.postgres.Close"

	if err := s.db.Close(); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
