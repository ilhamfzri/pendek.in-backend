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
	SQL_CHECK_USERNAME := `SELECT * 
							FROM users
							WHERE username=$1 LIMIT 1`
	rows, err := tx.QueryContext(ctx, SQL_CHECK_USERNAME, user.Username)
	helper.PanicIfError(err)
	if rows.Next() {
		return user, errors.New("username is used")
	}

	// check email is already registered
	SQL_CHECK_EMAIL := `SELECT * 
						FROM users 
						WHERE email=$1 LIMIT 1`
	rows, err = tx.QueryContext(ctx, SQL_CHECK_EMAIL, user.Email)
	helper.PanicIfError(err)
	if rows.Next() {
		return user, errors.New("email is already registered")
	}

	// insert new account to database
	SQL_CREATE_ACCOUNT := `INSERT INTO users (username, email, password, last_login, created_at, updated_at, verified) 
							VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`
	timeNow := time.Now().UTC()
	user.LastLogin = timeNow
	user.UpdatedAt = timeNow
	user.CreatedAt = timeNow

	err = tx.QueryRowContext(ctx, SQL_CREATE_ACCOUNT,
		user.Username,
		user.Email,
		user.Password,
		user.LastLogin,
		user.CreatedAt,
		user.UpdatedAt,
		false).Scan(&user.Id)

	helper.PanicIfError(err)
	return user, nil
}

func (repository *UserRepositoryImpl) CreateVerifyCode(ctx context.Context, tx *sql.Tx, user_id int, code string) error {
	SQL_DEACTIVATE_CURRENT_VERIFY := "UPDATE verify_and_forget_password SET active = false WHERE user_id = $1 AND type='VERIFY'"
	_, err := tx.ExecContext(ctx, SQL_DEACTIVATE_CURRENT_VERIFY, user_id)
	helper.PanicIfError(err)

	timeNow := time.Now().UTC()
	SQL_SET_VERIFY_CODE := "INSERT INTO verify_and_forget_password (user_id, type, code, active, date) VALUES ($1, $2, $3, $4, $5)"
	_, err = tx.ExecContext(ctx, SQL_SET_VERIFY_CODE, user_id, "VERIFY", code, true, timeNow)
	helper.PanicIfError(err)
	return err
}

func (repository *UserRepositoryImpl) Login(ctx context.Context, tx *sql.Tx, user domain.User) error {
	// check if email is exist
	SQL_CHECK_EMAIL := `SELECT verified,password 
						FROM users 
						WHERE email=$1 LIMIT 1`
	rows, err := tx.QueryContext(ctx, SQL_CHECK_EMAIL, user.Email)
	helper.PanicIfError(err)
	if !rows.Next() {
		return errors.New("email isn' registered")
	}

	var currentPassword string
	var verified bool
	err = rows.Scan(&verified, &currentPassword)
	helper.PanicIfError(err)

	// check if account is verified
	if !verified {
		return errors.New("account isn't verified")
	}

	// check if password is correct
	if currentPassword != user.Password {
		return errors.New("password incorrect")
	}
	return nil
}

func (repository *UserRepositoryImpl) Verify(ctx context.Context, tx *sql.Tx, email string, code string) error {
	SQL_VERIFY := `SELECT v.id
					FROM verify_and_forget_password as v
					LEFT JOIN users as u
					ON v.user_id=u.id
					WHERE 
						u.email=$1 AND
						v.code=$2 AND
						v.active=true
					LIMIT 1`
	rows, err := tx.QueryContext(ctx, SQL_VERIFY, email, code)
	helper.PanicIfError(err)

	if !rows.Next() {
		return errors.New("verification link isn't valid")
	}
	rows.Close()

	var verifyId int
	rows.Scan(&verifyId)

	SQL_UPDATE_VERIFY_STATUS_USER := `UPDATE users
									SET verified=true
									WHERE email=$1`
	_, err = tx.ExecContext(ctx, SQL_UPDATE_VERIFY_STATUS_USER, email)
	helper.PanicIfError(err)

	SQL_UPDATE_VERIFY_STATUS_ENTRY := `UPDATE verify_and_forget_password
											SET active=false
											WHERE id=$1`
	_, err = tx.ExecContext(ctx, SQL_UPDATE_VERIFY_STATUS_ENTRY, verifyId)
	helper.PanicIfError(err)

	return nil
}

func (repository *UserRepositoryImpl) Update(ctx context.Context, tx *sql.Tx, user domain.User) (domain.User, error) {
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
