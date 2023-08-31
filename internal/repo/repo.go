package repo

import (
	"context"

	"github.com/ABDURAZZAKK/avito_experiment/internal/entity"
	"github.com/ABDURAZZAKK/avito_experiment/internal/repo/pgdb"
	"github.com/ABDURAZZAKK/avito_experiment/pkg/postgres"
)

type User interface {
	Create(ctx context.Context, slug string) (int, error)
	GetById(ctx context.Context, id int) (entity.User, error)
	GetRandomIDs(ctx context.Context, limit int) ([]int, error)
	GetCount(ctx context.Context) (int, error)
	Delete(ctx context.Context, id int) (int, error)
}

type Segment interface {
	GetBySlug(ctx context.Context, slug string) (entity.Segment, error)
	Create(ctx context.Context, slug string) (string, error)
	CreateAll(ctx context.Context, slugs []string) error
	Delete(ctx context.Context, slug string) (string, error)
}

type UsersSegments interface {
	AddAndRemoveSegmentsUser(ctx context.Context, users []int, addList []string, removeList []string) error
	GetUserSegments(ctx context.Context, id int) ([]string, error)
	GetStatsPerPeriod(ctx context.Context, year int, month int) ([]entity.UsersSegmentsStats, error)
	DeleteSegmentFromUser(ctx context.Context, users []int, segments []string) error
}

type Repositories struct {
	User
	Segment
	UsersSegments
}

func NewRepositories(pg *postgres.Postgres) *Repositories {
	return &Repositories{
		User:          pgdb.NewUserRepo(pg),
		Segment:       pgdb.NewSegmentRepo(pg),
		UsersSegments: pgdb.NewUsersSegmentsRepo(pg),
	}
}
