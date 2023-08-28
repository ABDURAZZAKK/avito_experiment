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

type SegmentRepo struct {
	*postgres.Postgres
}

func NewSegmentRepo(pg *postgres.Postgres) *SegmentRepo {
	return &SegmentRepo{pg}
}

func (r *SegmentRepo) GetTransaction(ctx context.Context) (pgx.Tx, error) {
	tx, err := r.Pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("UsersSegmentsRepo.AddAndRemoveSegmentsUser - r.Pool.Begin: %v", err)
	}

	return tx, nil
}

func (r *SegmentRepo) GetBySlug(ctx context.Context, slug string) (entity.Segment, error) {
	sql, args, _ := r.Builder.
		Select("*").
		From("segments").
		Where("slug = ?", slug).
		ToSql()

	var segment entity.Segment
	err := r.Pool.QueryRow(ctx, sql, args...).Scan(
		&segment.Slug,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.Segment{}, repoerrs.ErrNotFound
		}
		return entity.Segment{}, fmt.Errorf("SegmentRepo.GetBySlug - r.Pool.QueryRow: %v", err)
	}

	return segment, nil
}

func (r *SegmentRepo) Create(ctx context.Context, slug string) (string, error) {
	sql, args, _ := r.Builder.
		Insert("segments").
		Columns("slug").
		Values(slug).
		Suffix("RETURNING slug").
		ToSql()

	var _slug string
	err := r.Pool.QueryRow(ctx, sql, args...).Scan(&_slug)
	if err != nil {
		log.Debugf("err: %v", err)
		var pgErr *pgconn.PgError
		if ok := errors.As(err, &pgErr); ok {
			if pgErr.Code == "23505" {
				return "", repoerrs.ErrAlreadyExists
			}
		}
		return "", fmt.Errorf("SegmentRepo.Create - r.Pool.QueryRow: %v", err)
	}
	return _slug, nil
}

func (r *SegmentRepo) CreateAll(ctx context.Context, slugs []string) error {
	builder := r.Builder.
		Insert("segments").
		Columns("slug")
	for _, slug := range slugs {
		builder = builder.Values(slug)
	}
	sql, args, _ := builder.ToSql()

	if _, err := r.Pool.Exec(ctx, sql, args...); err != nil {
		log.Debugf("err: %v", err)
		var pgErr *pgconn.PgError
		if ok := errors.As(err, &pgErr); ok {
			if pgErr.Code == "23505" {
				return repoerrs.ErrAlreadyExists
			}
		}
		return fmt.Errorf("SegmentRepo.CreateAll - r.Pool.QueryRow: %v", err)
	}
	return nil
}

func (r *SegmentRepo) Delete(ctx context.Context, slug string) (string, error) {
	sql, args, _ := r.Builder.
		Delete("segments").
		Where("slug = ?", slug).
		Suffix("RETURNING slug").
		ToSql()

	var s string
	err := r.Pool.QueryRow(ctx, sql, args...).Scan(&s)
	if err != nil {
		return "", fmt.Errorf("SegmentRepo.Delete - r.Pool.QueryRow: %v", err)
	}
	return s, nil
}
