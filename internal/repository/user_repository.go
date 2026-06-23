package repository

import (
	"database/sql"

	"github.com/google/uuid"
	"github.com/tetbatista/govault/internal/domain"
)

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) domain.UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(user *domain.User) error {
	user.ID = uuid.New().String()
	_, err := r.db.Exec(
		"INSERT INTO users (id, username, password, salt) VALUES (?, ?, ?, ?)",
		user.ID, user.Username, user.Password, user.Salt,
	)
	return err
}

func (r *userRepository) FindByUsername(username string) (*domain.User, error) {
	user := &domain.User{}
	err := r.db.QueryRow(
		"SELECT id, username, password, salt, created_at FROM users WHERE username = ?",
		username,
	).Scan(&user.ID, &user.Username, &user.Password, &user.Salt, &user.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return user, nil
}
