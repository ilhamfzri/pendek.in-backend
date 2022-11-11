package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/ilhamfzri/pendek.in/internal/helper"
	"github.com/ilhamfzri/pendek.in/internal/model/domain"
)

type UserRepositoryImpl struct {
}

func NewUserRepository() UserRepository {
	return &UserRepositoryImpl{}
}

func (repository *UserRepositoryImpl) Create(ctx context.Context, tx *sql.Tx, user domain.User) (domain.User, error) {

	// check username is already used
	SQL_CHECK_USERNAME := "SELECT * FROM users WHERE username=$1 LIMIT 1"
	rows, err := tx.QueryContext(ctx, SQL_CHECK_USERNAME, user.Username)
	helper.PanicIfError(err)
	if rows.Next() {
		return user, errors.New("username is already used")
	}

	// check email is already registered
	SQL_CHECK_EMAIL := "SELECT * FROM users WHERE email=$1 LIMIT 1"
	rows, err = tx.QueryContext(ctx, SQL_CHECK_EMAIL, user.Email)
	helper.PanicIfError(err)
	if rows.Next() {
		return user, errors.New("email is already registered")
	}

	// insert new account to database
	SQL_CREATE_ACCOUNT := "INSERT INTO users (first_name, last_name, username, email, password, activated_at, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)"
	timeNow := time.Now().UTC()
	user.ActivatedAt = timeNow
	user.UpdatedAt = timeNow
	user.CreatedAt = timeNow

	_, err = tx.ExecContext(ctx, SQL_CREATE_ACCOUNT, user.FirstName, user.LastName, user.Username, user.Email, user.Password, user.ActivatedAt, user.CreatedAt, user.UpdatedAt)
	helper.PanicIfError(err)

	// get last id
	rows, err = tx.QueryContext(ctx, "SELECT id FROM users WHERE email=$1", user.Email)
	helper.PanicIfError(err)

	var id int
	rows.Scan(&id)

	user.Id = int(id)
	return user, nil
}

func (repository *UserRepositoryImpl) Update(ctx context.Context, tx *sql.Tx, user domain.User) (domain.User, error) {
	panic("not implemented") // TODO: Implement
}

func (repository *UserRepositoryImpl) Login(ctx context.Context, tx *sql.Tx, user domain.User) (domain.User, error) {
	panic("not implemented") // TODO: Implement
}

func (repository *UserRepositoryImpl) FindByUsername(ctx context.Context, tx *sql.Tx, username string) (domain.User, error) {
	// check username is already used
	user := domain.User{}
	SQL_CHECK_USERNAME := "SELECT id, first_name, last_name, username, email FROM users WHERE username=$1 LIMIT 1"
	rows, err := tx.QueryContext(ctx, SQL_CHECK_USERNAME, username)
	helper.PanicIfError(err)
	if rows.Next() {
		err = rows.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Username, &user.Email)
		helper.PanicIfError(err)
		return user, nil
	} else {
		return user, errors.New("username is already used")
	}
}
