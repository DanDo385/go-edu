package database

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/user/go-edu/minis/50-mini-service-all-features/internal/models"
)

// DB provides in-memory database for demonstration
// In production, this would use sql.DB with PostgreSQL, MySQL, etc.
type DB struct {
	mu    sync.RWMutex
	users map[int]*models.User
	creds map[string]string // username -> password
}

// New creates a new in-memory database with sample data
func New() *DB {
	db := &DB{
		users: make(map[int]*models.User),
		creds: make(map[string]string),
	}

	// Add sample users
	db.users[1] = &models.User{
		ID:        1,
		Username:  "alice",
		Email:     "alice@example.com",
		CreatedAt: time.Now().Add(-30 * 24 * time.Hour),
	}
	db.users[2] = &models.User{
		ID:        2,
		Username:  "bob",
		Email:     "bob@example.com",
		CreatedAt: time.Now().Add(-15 * 24 * time.Hour),
	}
	db.users[3] = &models.User{
		ID:        3,
		Username:  "charlie",
		Email:     "charlie@example.com",
		CreatedAt: time.Now().Add(-7 * 24 * time.Hour),
	}

	// Add credentials (in production, use bcrypt!)
	db.creds["alice"] = "password123"
	db.creds["bob"] = "password123"
	db.creds["charlie"] = "password123"

	return db
}

// Health checks database connectivity
func (db *DB) Health(ctx context.Context) error {
	// In production, this would ping the actual database
	// For in-memory, we just check if the map is initialized
	db.mu.RLock()
	defer db.mu.RUnlock()

	if db.users == nil {
		return fmt.Errorf("database not initialized")
	}

	return nil
}

// Close closes database connections
func (db *DB) Close() error {
	// In production, this would close sql.DB
	// For in-memory, nothing to close
	return nil
}

// GetUser retrieves a user by ID
func (db *DB) GetUser(ctx context.Context, id int) (*models.User, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	user, ok := db.users[id]
	if !ok {
		return nil, fmt.Errorf("user not found")
	}

	return user, nil
}

// ListUsers returns all users
func (db *DB) ListUsers(ctx context.Context) ([]models.User, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	users := make([]models.User, 0, len(db.users))
	for _, user := range db.users {
		users = append(users, *user)
	}

	return users, nil
}

// Authenticate verifies username and password
func (db *DB) Authenticate(ctx context.Context, username, password string) (*models.User, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	// Check credentials
	storedPassword, ok := db.creds[username]
	if !ok || storedPassword != password {
		return nil, fmt.Errorf("invalid credentials")
	}

	// Find user by username
	for _, user := range db.users {
		if user.Username == username {
			return user, nil
		}
	}

	return nil, fmt.Errorf("user not found")
}
