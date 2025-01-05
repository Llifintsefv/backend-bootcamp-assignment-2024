package auth

import (
	"backend-bootcamp-assignment-2024/dto"
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

type AuthRepository interface {
	RegisterUser(ctx context.Context, req dto.PostRegisterJSONRequestBody, uuid string) error
	loginUser(ctx context.Context, req dto.PostLoginJSONRequestBody) (string, error)
}

type authRepository struct {
	db *sql.DB
}

func NewAuthRepository(db *sql.DB) AuthRepository {
	return &authRepository{
		db: db,
	}
}

func (r *authRepository) RegisterUser(ctx context.Context, req dto.PostRegisterJSONRequestBody, uuid string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(string(*req.Password)), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	stmt, err := r.db.PrepareContext(ctx, "INSERT INTO users (id, email, password, type) VALUES ($1, $2, $3, $4)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, uuid, req.Email, hashedPassword, req.UserType)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			if pqErr.Constraint == "users_email_key" {
				return errors.New("user with this email already exists")
			}
		}
		return err
	}

	return nil
}

func (r *authRepository) loginUser(ctx context.Context, req dto.PostLoginJSONRequestBody) (string, error) {
	stmt, err := r.db.PrepareContext(ctx, "SELECT password, type FROM users WHERE id = $1")
	if err != nil {
		return "", err
	}
	defer stmt.Close()

	var hashedPassword []byte
	var userType string
	err = stmt.QueryRowContext(ctx, req.Id).Scan(&hashedPassword, &userType)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("user not found")
		}
		return "", err
	}

	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(string(*req.Password)))
	if err != nil {

		return "", fmt.Errorf("invalid password")
	}

	return userType, nil
}
