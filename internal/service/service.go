package service

import (
	"context"

	"github.com/ABDURAZZAKK/avito_experiment/internal/entity"
	"github.com/ABDURAZZAKK/avito_experiment/internal/repo"
)

type User interface {
	Create(ctx context.Context, slug string) (int, error)
	GetRandomIDs(ctx context.Context, limit int) ([]int, error)
	GetRandomIDsWithPercent(ctx context.Context, percent int) ([]int, error)
	GetCount(ctx context.Context) (int, error)
	GetById(ctx context.Context, id int) (entity.User, error)
	ChangeSegments(ctx context.Context, id int, addList []string, removeList []string) error
	GetSegments(ctx context.Context, id int) ([]string, error)
	Delete(ctx context.Context, id int) (int, error)
}

type Segment interface {
	GetBySlug(ctx context.Context, slug string) (entity.Segment, error)
	Create(ctx context.Context, slug string, users []int) (string, error)
	CreateAll(ctx context.Context, slugs []string, users []int) error
	Delete(ctx context.Context, slug string) (string, error)
}

type Services struct {
	User
	Segment
}

type ServicesDependencies struct {
	Repos *repo.Repositories
}

func NewServices(deps ServicesDependencies) *Services {
	return &Services{
		User:    NewUserService(deps.Repos.User, deps.Repos.UsersSegments),
		Segment: NewSegmentService(deps.Repos.Segment, deps.Repos.UsersSegments, deps.Repos.User),
	}
}
