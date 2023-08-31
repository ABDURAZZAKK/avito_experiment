package pgdb

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/ABDURAZZAKK/avito_experiment/internal/entity"
	"github.com/ABDURAZZAKK/avito_experiment/internal/repo/repoerrs"
	"github.com/ABDURAZZAKK/avito_experiment/pkg/postgres"
	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type UsersSegmentsRepo struct {
	*postgres.Postgres
}

func NewUsersSegmentsRepo(pg *postgres.Postgres) *UsersSegmentsRepo {
	return &UsersSegmentsRepo{pg}
}

func (r *UsersSegmentsRepo) GetTransaction(ctx context.Context) (pgx.Tx, error) {
	tx, err := r.Pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("UsersSegmentsRepo.AddAndRemoveSegmentsUser - r.Pool.Begin: %v", err)
	}

	return tx, nil
}

func (r *UsersSegmentsRepo) getInsertSqlAddSegmentsToUser(users []int, segments []string, operation entity.Operation) (string, []interface{}, error) {
	builder := r.Builder.
		Insert("users_segments_stats").
		Columns("user_pk", "segment_pk", "created_at", "operation")
	for _, user := range users {
		for _, segment := range segments {
			builder = builder.
				Values(user, segment, time.Now(), operation)
		}
	}
	return builder.ToSql()
}
func (r *UsersSegmentsRepo) getDeleteUsersSegmentsSql(users []int, segments []string) (string, []interface{}, error) {
	some := squirrel.Or{}
	for _, user := range users {
		for _, slug := range segments {
			some = append(some, squirrel.Eq{"user_pk": user, "segment_pk": slug})
		}
	}
	return r.Builder.
		Delete("users_segments").
		Where(some).
		ToSql()
}

func getMonthStartEndDates(month int) (start string, end string) {
	_month := time.Month(month)
	start = time.Date(2023, _month, 1, 0, 0, 0, 0, time.Local).Format("2006-01-02")
	end = time.Date(2023, _month+1, 1, 0, 0, 0, 0, time.Local).Format("2006-01-02")
	return start, end
}

func (r *UsersSegmentsRepo) AddAndRemoveSegmentsUser(
	ctx context.Context,
	users []int,
	addList []string,
	removeList []string) error {
	tx, err := r.Pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("UsersSegmentsRepo.AddAndRemoveSegmentsUser - r.Pool.Begin: %v", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	if len(addList) != 0 {
		sql, args, _ := r.getInsertSqlAddSegmentsToUser(users, addList, entity.SEGMENT_ADDED)
		if _, err = tx.Exec(ctx, sql, args...); err != nil {
			return fmt.Errorf("UsersSegmentsRepo.AddAndRemoveSegmentsUser (add to stats) - tx.Exec: %v", err)
		}
		builder := r.Builder.
			Insert("users_segments").
			Columns("user_pk", "segment_pk")
		for _, user := range users {
			for _, segment := range addList {
				builder = builder.
					Values(user, segment)
			}
		}
		sql, args, _ = builder.ToSql()
		if _, err = tx.Exec(ctx, sql, args...); err != nil {
			var pgErr *pgconn.PgError
			if ok := errors.As(err, &pgErr); ok {
				if pgErr.Code == "23505" {
					return repoerrs.ErrAlreadyExists
				}
				if pgErr.Code == "23503" {
					return repoerrs.ErrNotFound
				}
			}
			return fmt.Errorf("UsersSegmentsRepo.AddAndRemoveSegmentsUser (add) - tx.Exec: %v", err)
		}

	}
	if len(removeList) != 0 {
		sql, args, _ := r.getInsertSqlAddSegmentsToUser(users, removeList, entity.SEGMENT_REMOVED)
		if _, err = tx.Exec(ctx, sql, args...); err != nil {
			return fmt.Errorf("UsersSegmentsRepo.AddAndRemoveSegmentsUser (remove to stats) - tx.Exec: %v", err)
		}

		sql, args, _ = r.getDeleteUsersSegmentsSql(users, removeList)

		if _, err = tx.Exec(ctx, sql, args...); err != nil {
			return fmt.Errorf("UsersSegmentsRepo.AddAndRemoveSegmentsUser (remove) - tx.Exec: %v", err)
		}
	}

	if err = tx.Commit(ctx); err != nil {
		return fmt.Errorf("UsersSegmentsRepo.AddAndRemoveSegmentsUser - tx.Commit: %v", err)
	}

	return nil
}

func (r *UsersSegmentsRepo) GetUserSegments(ctx context.Context, id int) ([]string, error) {
	sql, args, _ := r.Builder.
		Select("segment_pk").
		From("users_segments").
		Where("user_pk = ?", id).
		ToSql()

	rows, err := r.Pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("UsersSegmentsRepo.GetUserSegments - r.Pool.Query: %v", err)
	}
	defer rows.Close()
	var segments []string
	for rows.Next() {
		var s string
		err := rows.Scan(&s)
		if err != nil {
			return nil, fmt.Errorf("UsersSegmentsRepo.GetUserSegments - rows.Scan: %v", err)
		}
		segments = append(segments, s)
	}

	if err = rows.Err(); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, repoerrs.ErrNotFound
		}
		return nil, fmt.Errorf("UsersSegmentsRepo.GetUserSegments -rows.Err: %v", err)
	}

	return segments, nil
}

func (r *UsersSegmentsRepo) GetStatsPerPeriod(ctx context.Context, year int, month int) ([]entity.UsersSegmentsStats, error) {
	month_start, month_end := getMonthStartEndDates(month)
	sql, args, _ := r.Builder.
		Select("*").
		From("users_segments_stats").
		Where("created_at >= ? AND created_at < ?", month_start, month_end).
		ToSql()

	rows, err := r.Pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("UsersSegmentsRepo.GetUserSegments - r.Pool.Query: %v", err)
	}
	defer rows.Close()

	var segments []entity.UsersSegmentsStats
	for rows.Next() {
		var s entity.UsersSegmentsStats
		err := rows.Scan(
			&s.User,
			&s.Segment,
			&s.Created_at,
			&s.Operation,
		)
		if err != nil {
			return nil, fmt.Errorf("UsersSegmentsRepo.GetUserSegments - rows.Scan: %v", err)
		}
		segments = append(segments, s)
	}

	if err = rows.Err(); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, repoerrs.ErrNotFound
		}
		return nil, fmt.Errorf("UsersSegmentsRepo.GetUserSegments -rows.Err: %v", err)
	}

	return segments, nil

}
func (r *UsersSegmentsRepo) DeleteSegmentFromUser(ctx context.Context, users []int, segments []string) error {
	tx, err := r.Pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("UsersSegmentsRepo.DeleteSegmentFromUser - r.Pool.Begin: %v", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	if len(users) != 0 && len(segments) != 0 {
		sql, args, _ := r.getInsertSqlAddSegmentsToUser(users, segments, entity.SEGMENT_REMOVED)
		if _, err = tx.Exec(ctx, sql, args...); err != nil {
			return fmt.Errorf("UsersSegmentsRepo.DeleteSegmentFromUser (remove to stats) - tx.Exec: %v", err)
		}

		sql, args, _ = r.getDeleteUsersSegmentsSql(users, segments)

		if _, err = tx.Exec(ctx, sql, args...); err != nil {
			return fmt.Errorf("UsersSegmentsRepo.DeleteSegmentFromUser (remove) - tx.Exec: %v", err)
		}
	}
	if err = tx.Commit(ctx); err != nil {
		return fmt.Errorf("UsersSegmentsRepo.DeleteSegmentFromUser - tx.Commit: %v", err)
	}

	return nil
}
