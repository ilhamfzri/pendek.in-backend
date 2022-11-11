package repository

import (
	"context"
	"database/sql"

	"github.com/ilhamfzri/pendek.in/internal/model/domain"
)

type UserRepository interface {
	Create(ctx context.Context, tx *sql.Tx, user domain.User) (domain.User, error)
	Update(ctx context.Context, tx *sql.Tx, user domain.User) (domain.User, error)
	Login(ctx context.Context, tx *sql.Tx, user domain.User) (domain.User, error)
	FindByUsername(ctx context.Context, tx *sql.Tx, username string) (domain.User, error)
}
