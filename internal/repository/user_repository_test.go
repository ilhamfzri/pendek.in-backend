package repository

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/ilhamfzri/pendek.in/internal/helper"
	"github.com/ilhamfzri/pendek.in/internal/model/domain"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func initTestUserDB() *sql.DB {
	connStr := "postgres://postgres:example@localhost:5432/pendekin_db_test?sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	helper.PanicIfError(err)

	err = db.Ping()
	helper.PanicIfError(err)

	db.SetMaxIdleConns(5)
	db.SetMaxOpenConns(20)
	db.SetConnMaxLifetime(60 * time.Minute)

	clearUserTable(db)
	return db
}

func clearUserTable(db *sql.DB) {
	ctx := context.Background()
	_, err := db.ExecContext(ctx, "DELETE FROM verify_and_forget_password WHERE id > 0")
	helper.PanicIfError(err)
	_, err = db.ExecContext(ctx, "DELETE FROM users WHERE id > 1")
	helper.PanicIfError(err)
}

func insertDummyUser(db *sql.DB, userRepository UserRepository) (int, string) {
	ctx := context.Background()
	tx, err := db.Begin()
	helper.PanicIfError(err)
	defer helper.CommitOrRollback(tx)

	user := domain.User{
		Username: "usertest",
		Email:    "usertest@example.com",
		Password: "example_hash_password", //Hashed Password
	}
	user, err = userRepository.Create(ctx, tx, user)
	helper.PanicIfError(err)
	return user.Id, user.Email
}

func insertDummyVerifiedUser(db *sql.DB, userRepository UserRepository) {
	ctx := context.Background()
	tx, err := db.Begin()
	helper.PanicIfError(err)
	defer helper.CommitOrRollback(tx)

	user := domain.User{
		Username: "usertest",
		Email:    "usertest@example.com",
		Password: "example_hash_password", //Hashed Password
	}
	user, err = userRepository.Create(ctx, tx, user)
	helper.PanicIfError(err)
	_, err = tx.ExecContext(ctx, `UPDATE users
								SET verified=true
								WHERE email=$1`, user.Email)
	helper.PanicIfError(err)
}

func insertDummyVerifyData(db *sql.DB, userRepository UserRepository, id int) {
	ctx := context.Background()
	tx, err := db.Begin()
	helper.PanicIfError(err)
	defer helper.CommitOrRollback(tx)

	err = userRepository.CreateVerifyCode(ctx, tx, id, "TESTCODE")
	helper.PanicIfError(err)
}

func TestUserRepositoryCreate(t *testing.T) {
	userRepository := NewUserRepository()
	db := initTestUserDB()

	t.Run("[User Repo][Create][Success]", func(t *testing.T) {
		userCreateSuccess(t, userRepository, db)
		clearUserTable(db)
	})

	t.Run("[User Repo][Create][Failed:Username Used]", func(t *testing.T) {
		insertDummyUser(db, userRepository)
		userCreateFailedUsernameUsed(t, userRepository, db)
		clearUserTable(db)
	})

	t.Run("[User Repo][Create][Failed:Email Registered]", func(t *testing.T) {
		insertDummyUser(db, userRepository)
		userCreateFailedEmailRegistered(t, userRepository, db)
		clearUserTable(db)
	})
}

func userCreateSuccess(t *testing.T, userRepository UserRepository, db *sql.DB) {
	ctx := context.Background()
	tx, err := db.Begin()
	defer helper.CommitOrRollback(tx)
	helper.PanicIfError(err)

	user := domain.User{
		Username: "usertest",
		Email:    "usertest@example.com",
		Password: "example_hash_password", //Hashed Password
	}
	user, err = userRepository.Create(ctx, tx, user)

	assert.Nil(t, err)
	assert.Equal(t, "usertest", user.Username)
	assert.Equal(t, "usertest@example.com", user.Email)
}

func userCreateFailedUsernameUsed(t *testing.T, userRepository UserRepository, db *sql.DB) {
	ctx := context.Background()
	tx, err := db.Begin()
	helper.PanicIfError(err)
	user := domain.User{
		Username: "usertest",
		Email:    "usertes2t@example.com",
		Password: "example_hash_password", //Hashed Password
	}
	user, err = userRepository.Create(ctx, tx, user)
	helper.CommitOrRollback(tx)

	assert.NotNil(t, err)
	assert.Equal(t, "username is used", err.Error())
}

func userCreateFailedEmailRegistered(t *testing.T, userRepository UserRepository, db *sql.DB) {
	ctx := context.Background()
	tx, err := db.Begin()
	helper.PanicIfError(err)

	user := domain.User{
		Username: "usertest2",
		Email:    "usertest@example.com",
		Password: "example_hash_password", //Hashed Password
	}
	user, err = userRepository.Create(ctx, tx, user)
	helper.CommitOrRollback(tx)

	assert.NotNil(t, err)
	assert.Equal(t, "email is already registered", err.Error())
}

func TestUserRepositoryVerify(t *testing.T) {
	userRepository := NewUserRepository()
	db := initTestUserDB()

	t.Run("[User Repo][CreateVerifyCode][Success]", func(t *testing.T) {
		id, _ := insertDummyUser(db, userRepository)
		userCreateVerifyCode(t, userRepository, db, id)
		clearUserTable(db)
	})

	t.Run("[User Repo][Verify][Success]", func(t *testing.T) {
		id, email := insertDummyUser(db, userRepository)
		insertDummyVerifyData(db, userRepository, id)
		userVerifyCodeSuccess(t, userRepository, db, email)
		clearUserTable(db)
	})

	t.Run("[User Repo][Verify][Failed]", func(t *testing.T) {
		id, email := insertDummyUser(db, userRepository)
		insertDummyVerifyData(db, userRepository, id)
		userVerifyCodeFailed(t, userRepository, db, email)
		clearUserTable(db)
	})

}

func userCreateVerifyCode(t *testing.T, userRepository UserRepository, db *sql.DB, id int) {
	ctx := context.Background()
	tx, err := db.Begin()
	helper.PanicIfError(err)
	defer helper.CommitOrRollback(tx)

	err = userRepository.CreateVerifyCode(ctx, tx, id, "TESTCODE")
	assert.Nil(t, err)
}

func userVerifyCodeSuccess(t *testing.T, userRepository UserRepository, db *sql.DB, email string) {
	ctx := context.Background()
	tx, err := db.Begin()
	helper.PanicIfError(err)
	defer helper.CommitOrRollback(tx)

	err = userRepository.Verify(ctx, tx, email, "TESTCODE")
	assert.Nil(t, err)
}

func userVerifyCodeFailed(t *testing.T, userRepository UserRepository, db *sql.DB, email string) {
	ctx := context.Background()
	tx, err := db.Begin()
	helper.PanicIfError(err)
	defer helper.CommitOrRollback(tx)

	err = userRepository.Verify(ctx, tx, email, "NOTVALID")
	assert.NotNil(t, err)
	assert.Equal(t, "verification link isn't valid", err.Error())
}

func TestUserRepositoryLogin(t *testing.T) {
	userRepository := NewUserRepository()
	db := initTestUserDB()

	t.Run("[User Repo][Login][Success]", func(t *testing.T) {
		insertDummyVerifiedUser(db, userRepository)
		userLoginSuccess(t, userRepository, db)
		clearUserTable(db)
	})

	t.Run("[User Repo][Login][Failed:Email Isn't Registered]", func(t *testing.T) {
		insertDummyVerifiedUser(db, userRepository)
		userLoginFailedEmailNotRegistered(t, userRepository, db)
		clearUserTable(db)
	})

	t.Run("[User Repo][Login][Failed:Password Incorrect]", func(t *testing.T) {
		insertDummyVerifiedUser(db, userRepository)
		userLoginFailedPasswordIncorrect(t, userRepository, db)
		clearUserTable(db)
	})
	t.Run("[User Repo][Login][Failed:Account Isnt Verified]", func(t *testing.T) {
		insertDummyUser(db, userRepository)
		userLoginFailedUserNotVerified(t, userRepository, db)
		clearUserTable(db)
	})
}

func userLoginSuccess(t *testing.T, userRepository UserRepository, db *sql.DB) {
	ctx := context.Background()
	tx, err := db.Begin()
	helper.PanicIfError(err)
	defer helper.CommitOrRollback(tx)

	user := domain.User{
		Username: "usertest",
		Email:    "usertest@example.com",
		Password: "example_hash_password", //Hashed Password
	}
	err = userRepository.Login(ctx, tx, user)
	assert.Nil(t, err)
}

func userLoginFailedEmailNotRegistered(t *testing.T, userRepository UserRepository, db *sql.DB) {
	ctx := context.Background()
	tx, err := db.Begin()
	helper.PanicIfError(err)
	defer helper.CommitOrRollback(tx)

	user := domain.User{
		Username: "usertest",
		Email:    "usertest2@example.com",
		Password: "example_hash_password", //Hashed Password
	}
	err = userRepository.Login(ctx, tx, user)
	assert.NotNil(t, err)
	assert.Equal(t, "email isn' registered", err.Error())
}

func userLoginFailedPasswordIncorrect(t *testing.T, userRepository UserRepository, db *sql.DB) {
	ctx := context.Background()
	tx, err := db.Begin()
	helper.PanicIfError(err)
	defer helper.CommitOrRollback(tx)

	user := domain.User{
		Username: "usertest",
		Email:    "usertest@example.com",
		Password: "example_wrong_hash_password", //Hashed Password
	}
	err = userRepository.Login(ctx, tx, user)
	assert.NotNil(t, err)
	assert.Equal(t, "password incorrect", err.Error())
}

func userLoginFailedUserNotVerified(t *testing.T, userRepository UserRepository, db *sql.DB) {
	ctx := context.Background()
	tx, err := db.Begin()
	helper.PanicIfError(err)
	defer helper.CommitOrRollback(tx)

	user := domain.User{
		Username: "usertest",
		Email:    "usertest@example.com",
		Password: "example_hash_password", //Hashed Password
	}
	err = userRepository.Login(ctx, tx, user)
	assert.NotNil(t, err)
	assert.Equal(t, "account isn't verified", err.Error())
}

func TestUserRepositoryUpdate(t *testing.T) {
	userRepository := NewUserRepository()
	db := initTestUserDB()

	t.Run("[User Repo][Update][Success]", func(t *testing.T) {
		insertDummyVerifiedUser(db, userRepository)
		userUpdateSuccess(t, userRepository, db)
		clearUserTable(db)
	})

}

func userUpdateSuccess(t *testing.T, userRepository UserRepository, db *sql.DB) {
	ctx := context.Background()
	tx, err := db.Begin()
	helper.PanicIfError(err)
	defer helper.CommitOrRollback(tx)

	user := domain.User{
		Username:  "usertest",
		FirstName: "test_firstname",
		LastName:  "test_lastname",
		Bio:       "test_bio",
	}

	user, err = userRepository.Update(ctx, tx, user)
	assert.Equal(t, "usertest", user.Username)
	assert.Equal(t, "test_firstname", user.FirstName)
	assert.Equal(t, "test_lastname", user.LastName)
	assert.Equal(t, "test_bio", user.Bio)
	assert.Nil(t, err)
}
