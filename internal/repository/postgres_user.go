package repository

import (
	"context"
	"github.com/jackc/pgx/v5/pgtype"
	"go-api/internal/db"
)

type postgresUserRepo struct {
	q *db.Queries
}

func NewPostgresUserRepository(q *db.Queries) UserRepository {
	return &postgresUserRepo{
		q: q,
	}
}

func (r *postgresUserRepo) CreateUser(ctx context.Context, arg db.CreateUserParams) (db.User, error) {
	return r.q.CreateUser(ctx, arg)
}

func (r *postgresUserRepo) GetUserByEmail(ctx context.Context, email string) (db.User, error) {
	return r.q.GetUserByEmail(ctx, email)
}

func (r *postgresUserRepo) CreateSession(ctx context.Context, arg db.CreateSessionParams) (db.Session, error) {
	return r.q.CreateSession(ctx, arg)
}

func (r *postgresUserRepo) GetSession(ctx context.Context, id pgtype.UUID) (db.Session, error) {
	return r.q.GetSession(ctx, id)
}
