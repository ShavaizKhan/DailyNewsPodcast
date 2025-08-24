package graph

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	"github.com/ShavaizKhan/DailyNewsPodcast/webapp-backend/graph/model"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
	_ "github.com/joho/godotenv/autoload"
)

// UserStore defines the interface for our database operations.
type UserStore interface {
	CreateUser(ctx context.Context, email, passwordHash, country, topic string) (*model.User, error)
	GetUserByEmail(ctx context.Context, email string) (*model.User, error)
	GetUserByID(ctx context.Context, id string) (*model.User, error)
	UpdateUserPreferences(ctx context.Context, id, country, topic string) (*model.Preferences, error)
}

// PGStore implements UserStore using the standard library.
type PGStore struct {
	db *sql.DB
}

// NewPGStore creates a new PGStore and connects to the database.
func NewPGStore(ctx context.Context) (*PGStore, error) {
	godotenv.Load()
	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		return nil, fmt.Errorf("DATABASE_URL environment variable is not set")
	}

	// Use sql.Open with the "pgx" driver name.
	db, err := sql.Open("pgx", connStr)
	if err != nil {
		return nil, fmt.Errorf("unable to open database connection: %w", err)
	}

	// Ping the database to ensure the connection is live.
	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("could not ping database: %w", err)
	}

	return &PGStore{db: db}, nil
}

// Close closes the database connection pool.
func (p *PGStore) Close() {
	p.db.Close()
}

// CreateUser saves a new user to the database.
func (p *PGStore) CreateUser(ctx context.Context, email, passwordHash, country, topic string) (*model.User, error) {
	var user model.User
	var preferences model.Preferences
	query := `
		INSERT INTO users (email, password_hash, country, topic)
		VALUES ($1, $2, $3, $4)
		RETURNING id, email, country, topic`

	err := p.db.QueryRowContext(ctx, query, email, passwordHash, country, topic).Scan(&user.ID, &user.Email, &preferences.Country, &preferences.Topic)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	user.Preferences = &preferences
	return &user, nil
}

// GetUserByEmail fetches a user by their email address.
func (p *PGStore) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User
	var preferences model.Preferences
	query := `
		SELECT id, email, password_hash, country, topic
		FROM users
		WHERE email = $1`

	err := p.db.QueryRowContext(ctx, query, email).Scan(&user.ID, &user.Email, &user.PasswordHash, &preferences.Country, &preferences.Topic)
	if err != nil {
		return nil, fmt.Errorf("user with email '%s' not found: %w", email, err)
	}

	user.Preferences = &preferences
	return &user, nil
}

// GetUserByID fetches a user by their ID.
func (p *PGStore) GetUserByID(ctx context.Context, id string) (*model.User, error) {
	var user model.User
	var preferences model.Preferences
	query := `
		SELECT id, email, country, topic
		FROM users
		WHERE id = $1`

	err := p.db.QueryRowContext(ctx, query, id).Scan(&user.ID, &user.Email, &preferences.Country, &preferences.Topic)
	if err != nil {
		return nil, fmt.Errorf("user with ID '%s' not found: %w", id, err)
	}

	user.Preferences = &preferences
	return &user, nil
}

// UpdateUserPreferences updates a user's country and topic preferences.
func (p *PGStore) UpdateUserPreferences(ctx context.Context, id, country, topic string) (*model.Preferences, error) {
	var preferences model.Preferences
	query := `
		UPDATE users
		SET country = $1, topic = $2
		WHERE id = $3
		RETURNING country, topic`

	err := p.db.QueryRowContext(ctx, query, country, topic, id).Scan(&preferences.Country, &preferences.Topic)
	if err != nil {
		return nil, fmt.Errorf("failed to update user preferences: %w", err)
	}

	return &preferences, nil
}
