package repository

import (
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/tetbatista/govault/internal/domain"
)

type credentialRepository struct {
	db *sql.DB
}

func NewCredentialRepository(db *sql.DB) domain.CredentialRepository {
	return &credentialRepository{db: db}
}

func (r *credentialRepository) Create(credential *domain.Credential) error {
	credential.ID = uuid.New().String()
	_, err := r.db.Exec(
		"INSERT INTO credentials (id, user_id, site, username, password) VALUES (?, ?, ?, ?, ?)",
		credential.ID, credential.UserID, credential.Site, credential.Username, credential.Password,
	)
	return err
}

func (r *credentialRepository) FindAllByUserID(userID string) ([]domain.Credential, error) {
	rows, err := r.db.Query(
		"SELECT id, user_id, site, username, password, created_at FROM credentials WHERE user_id = ?",
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var credentials []domain.Credential
	for rows.Next() {
		var c domain.Credential
		err := rows.Scan(&c.ID, &c.UserID, &c.Site, &c.Username, &c.Password, &c.CreatedAt)
		if err != nil {
			return nil, err
		}
		credentials = append(credentials, c)
	}
	return credentials, nil
}

func (r *credentialRepository) FindByUserIDAndSite(userID string, site string) (*domain.Credential, error) {
	c := &domain.Credential{}
	err := r.db.QueryRow(
		"SELECT id, user_id, site, username, password, created_at FROM credentials WHERE user_id = ? AND site = ?",
		userID, site,
	).Scan(&c.ID, &c.UserID, &c.Site, &c.Username, &c.Password, &c.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (r *credentialRepository) DeleteByUserIDAndSite(userID string, site string) error {
	result, err := r.db.Exec(
		"DELETE FROM credentials WHERE user_id = ? AND site = ?",
		userID, site,
	)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return errors.New("credencial nao encontrada para o site: " + site)
	}
	return nil
}
