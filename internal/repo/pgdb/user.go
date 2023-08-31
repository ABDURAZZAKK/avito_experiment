package pgdb

import (
	"context"
	"errors"
	"fmt"

	"github.com/ABDURAZZAKK/avito_experiment/internal/entity"
	"github.com/ABDURAZZAKK/avito_experiment/internal/repo/repoerrs"
	"github.com/ABDURAZZAKK/avito_experiment/pkg/postgres"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	log "github.com/sirupsen/logrus"
)

type UserRepo struct {
	*postgres.Postgres
}

func NewUserRepo(pg *postgres.Postgres) *UserRepo {
	return &UserRepo{pg}
}

func (r *UserRepo) GetTransaction(ctx context.Context) (pgx.Tx, error) {
	tx, err := r.Pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("UsersSegmentsRepo.AddAndRemoveSegmentsUser - r.Pool.Begin: %v", err)
	}

	return tx, nil
}

func (r *UserRepo) Create(ctx context.Context, slug string) (int, error) {
	sql, args, _ := r.Builder.
		Insert("users").
		Columns("slug").
		Values(slug).
		Suffix("RETURNING id").
		ToSql()

	var id int
	err := r.Pool.QueryRow(ctx, sql, args...).Scan(&id)
	if err != nil {
		log.Debugf("err: %v", err)
		var pgErr *pgconn.PgError
		if ok := errors.As(err, &pgErr); ok {
			if pgErr.Code == "23505" {
				return 0, repoerrs.ErrAlreadyExists
			}
		}
		return 0, fmt.Errorf("UserRepo.Create - r.Pool.QueryRow: %v", err)
	}
	return id, nil
}

func (r *UserRepo) GetById(ctx context.Context, id int) (entity.User, error) {
	sql, args, _ := r.Builder.
		Select("*").
		From("users").
		Where("id = ?", id).
		ToSql()

	var user entity.User
	err := r.Pool.QueryRow(ctx, sql, args...).Scan(
		&user.Id,
		&user.Slug,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.User{}, repoerrs.ErrNotFound
		}
		return entity.User{}, fmt.Errorf("UserRepo.GetById - r.Pool.QueryRow: %v", err)
	}

	return user, nil
}

func (r *UserRepo) Delete(ctx context.Context, id int) (int, error) {
	sql, args, _ := r.Builder.
		Delete("users").
		Where("id = ?", id).
		Suffix("RETURNING id").
		ToSql()

	var u_id int
	err := r.Pool.QueryRow(ctx, sql, args...).Scan(&u_id)
	if err != nil {
		return 0, fmt.Errorf("UserRepo.Delete - r.Pool.QueryRow: %v", err)
	}
	return u_id, nil

}

func (r *UserRepo) GetRandomIDs(ctx context.Context, limit int) ([]int, error) {
	sql := fmt.Sprintf("SELECT id FROM users ORDER BY random() LIMIT %d;", limit)

	rows, err := r.Pool.Query(ctx, sql)
	if err != nil {
		return nil, fmt.Errorf("UserRepo.GetRandomIDs - r.Pool.Query: %v", err)
	}
	defer rows.Close()
	var users []int
	for rows.Next() {
		var i int
		err := rows.Scan(&i)
		if err != nil {
			return nil, fmt.Errorf("UserRepo.GetRandomIDs - rows.Scan: %v", err)
		}
		users = append(users, i)
	}

	if err = rows.Err(); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, repoerrs.ErrNotFound
		}
		return nil, fmt.Errorf("UserRepo.GetRandomIDs - rows.Err: %v", err)
	}

	return users, nil
}

func (r *UserRepo) GetCount(ctx context.Context) (int, error) {
	sql := "SELECT * FROM users"

	rows, err := r.Pool.Query(ctx, sql)
	if err != nil {
		return 0, fmt.Errorf("UserRepo.GetCount - r.Pool.Query: %v", err)
	}
	defer rows.Close()
	var count int
	for rows.Next() {
		count++
	}
	if err != nil {
		return 0, fmt.Errorf("UserRepo.GetCount - rows.Scan: %v", err)
	}
	return count, nil
}
