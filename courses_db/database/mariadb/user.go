package mariadb

import (
	"database/sql"

	"github.com/antoniofmoliveira/courses/dto"
	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	c := &UserRepository{
		db: db,
	}
	c.db.Exec("CREATE TABLE IF NOT EXISTS users (id CHAR(36) PRIMARY KEY, name TEXT, email TEXT, password TEXT)")
	return c
}

func (r *UserRepository) FindByEmail(email string) (*dto.GetJWTInput, error) {
	var password string
	err := r.db.QueryRow("SELECT password FROM users WHERE email = ?", email).
		Scan(&password)
	if err != nil {
		return nil, err
	}
	return &dto.GetJWTInput{Email: email, Password: password}, nil
}

func (r *UserRepository) Create(user dto.UserInputDto) (dto.UserOutputDto, error) {
	id := uuid.New().String()
	_, err := r.db.Exec("INSERT INTO users (id, name, email, password) VALUES (?,?, ?, ?)",
		id, user.Name, user.Email, user.Password)
	if err != nil {
		return dto.UserOutputDto{}, err
	}
	return dto.UserOutputDto{ID: id, Name: user.Name, Email: user.Email}, nil
}

func (r *UserRepository) Delete(id string) error {
	_, err := r.db.Exec("DELETE FROM users WHERE id = ?", id)
	if err != nil {
		return err
	}
	return nil
}

func (r *UserRepository) Update(user dto.UserInputDto) error {
	_, err := r.db.Exec("UPDATE users SET name = ?, email = ?, password = ? WHERE id = ?",
		user.Name, user.Email, user.Password, user.ID)
	if err != nil {
		return err
	}
	return nil
}

func (r *UserRepository) FindAll() (dto.UserListOutputDto, error) {
	rows, err := r.db.Query("SELECT id, name, email FROM users")
	if err != nil {
		return dto.UserListOutputDto{}, err
	}
	defer rows.Close()
	var users dto.UserListOutputDto
	for rows.Next() {
		var user dto.UserOutputDto
		if err := rows.Scan(&user.ID, &user.Name, &user.Email); err != nil {
			return dto.UserListOutputDto{}, err
		}
		users.Users = append(users.Users, user)
	}
	return users, nil
}

func (r *UserRepository) Find(id string) (dto.UserOutputDto, error) {
	var user dto.UserOutputDto
	err := r.db.QueryRow("SELECT id, name, email FROM users WHERE id = ?", id).
		Scan(&user.ID, &user.Name, &user.Email)
	if err != nil {
		return dto.UserOutputDto{}, err
	}
	return user, nil
}
