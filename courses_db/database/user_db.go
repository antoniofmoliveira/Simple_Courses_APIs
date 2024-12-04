package database

import (
	"database/sql"

	"github.com/antoniofmoliveira/courses/dto"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	c := &UserRepository{
		db: db,
	}
	c.db.Exec("CREATE TABLE IF NOT EXISTS users (id TEXT PRIMARY KEY, name TEXT, email TEXT, password TEXT)")
	return c
}

func (r *UserRepository) FindByEmail(email string) (*dto.GetJWTInput, error) {
	var password string
	err := r.db.QueryRow("SELECT password FROM users WHERE email = $1", email).
		Scan(&password)
	if err != nil {
		return nil, err
	}
	return &dto.GetJWTInput{Email: email, Password: password}, nil
}

func (r *UserRepository) Create(user *dto.CreateUserInput) error {
	_, err := r.db.Exec("INSERT INTO users (name, email, password) VALUES ($1, $2, $3)",
		user.Name, user.Email, user.Password)
	if err != nil {
		return err
	}
	return nil
}
