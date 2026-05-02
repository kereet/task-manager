package database

import (
	"database/sql"
	"fmt"
	"task-manager/internal/auth"
	"task-manager/internal/models"
	"time"

	"github.com/jmoiron/sqlx"
)

type UserStore struct {
	db *sqlx.DB
}

func NewUserStore(db *sqlx.DB) *UserStore {
	return &UserStore{db: db}
}

func (s *UserStore) Create(input models.RegisterInput) (*models.User, error) {
	var user models.User
	hashedPassword, err := auth.HashPassword(input.Password)
	if err != nil {
		return nil, err
	}
	query := `INSERT INTO users (name, email, password_hash, created_at) 
VALUES ($1, $2, $3, $4) 
RETURNING *;`
	now := time.Now()
	err = s.db.QueryRowx(query, input.Name, input.Email, hashedPassword, now).StructScan(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *UserStore) GetByEmail(email string) (*models.User, error) {
	var user models.User

	query := `SELECT * FROM users WHERE email = $1;`
	err := s.db.Get(&user, query, email)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user with email %s not found", email)
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *UserStore) GetByID(id int) (*models.User, error) {
	var user models.User

	query := `SELECT * FROM users WHERE id = $1`
	err := s.db.Get(&user, query, id)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user with id %d not found", id)
	}

	if err != nil {
		return nil, err
	}

	return &user, nil
}
