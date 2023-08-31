package service

import (
	"context"
	"fmt"

	"github.com/ABDURAZZAKK/avito_experiment/internal/entity"
	"github.com/ABDURAZZAKK/avito_experiment/internal/repo"
	"github.com/ABDURAZZAKK/avito_experiment/internal/repo/repoerrs"
)

type UserService struct {
	userRepo          repo.User
	usersSegmentsRepo repo.UsersSegments
}

func NewUserService(userRepo repo.User, usersSegmentsRepo repo.UsersSegments) *UserService {
	return &UserService{userRepo: userRepo, usersSegmentsRepo: usersSegmentsRepo}
}

func (s *UserService) Create(ctx context.Context, slug string) (int, error) {
	id, err := s.userRepo.Create(ctx, slug)
	if err != nil {
		if err == repoerrs.ErrAlreadyExists {
			return 0, ErrUserAlreadyExists
		}
		return 0, fmt.Errorf("UserService.Create - userRepo.Create: %v", err)
	}

	return id, nil
}

func (s *UserService) GetById(ctx context.Context, user_pk int) (entity.User, error) {
	user, err := s.userRepo.GetById(ctx, user_pk)
	if err != nil {
		if err == repoerrs.ErrNotFound {
			return entity.User{}, ErrUserNotFound
		}
		return entity.User{}, fmt.Errorf("UserService.GetById - userRepo.GetById: %v", err)
	}
	return user, nil
}

func (s *UserService) ChangeSegments(ctx context.Context, user_pk int, addList []string, removeList []string) error {
	_, err := s.userRepo.GetById(ctx, user_pk)
	if err != nil {
		if err == repoerrs.ErrNotFound {
			return ErrUserNotFound
		}
		return fmt.Errorf("UserService.ChangeSegments - userRepo.GetById: %v", err)
	}

	err = s.usersSegmentsRepo.AddAndRemoveSegmentsUser(ctx, []int{user_pk}, addList, removeList)
	if err != nil {
		if err == repoerrs.ErrAlreadyExists {
			return ErrAlreadyExists
		}
		if err == repoerrs.ErrNotFound {
			return ErrNotFound
		}
		return fmt.Errorf("UserService.ChangeSegments - usersSegmentsRepo.AddAndRemoveSegmentsUser: %v", err)
	}
	return nil
}

func (s *UserService) GetSegments(ctx context.Context, id int) ([]string, error) {
	_, err := s.userRepo.GetById(ctx, id)
	if err != nil {
		if err == repoerrs.ErrNotFound {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("UserService.GetSegments - userRepo.GetById: %v", err)
	}
	return s.usersSegmentsRepo.GetUserSegments(ctx, id)
}

func (s *UserService) Delete(ctx context.Context, id int) (int, error) {
	_, err := s.userRepo.GetById(ctx, id)
	if err != nil {
		if err == repoerrs.ErrNotFound {
			return 0, ErrUserNotFound
		}
		return 0, fmt.Errorf("UserService.Delete - userRepo.GetById: %v", err)
	}
	u_id, err := s.userRepo.Delete(ctx, id)
	if err != nil {
		return 0, fmt.Errorf("UserService.Delete - userRepo.Delete: %v", err)
	}
	return u_id, nil
}

func (s *UserService) GetRandomIDs(ctx context.Context, limit int) ([]int, error) {
	return s.userRepo.GetRandomIDs(ctx, limit)

}
func (s *UserService) GetCount(ctx context.Context) (int, error) {
	return s.userRepo.GetCount(ctx)
}

func (s *UserService) GetRandomIDsWithPercent(ctx context.Context, percent int) ([]int, error) {
	if percent > 0 {
		count, err := s.userRepo.GetCount(ctx)
		if err != nil {
			return nil, fmt.Errorf("UserService.GetRandomIDsWithPercent - userRepo.GetCount: %v", err)
		}

		limit := int(float32(count) * (float32(percent) / 100))
		if limit < 1 {
			limit++
		}
		users, err := s.userRepo.GetRandomIDs(ctx, limit)
		if err != nil {
			return nil, fmt.Errorf("UserService.GetRandomIDsWithPercent - userRepo.GetRandomIDs: %v", err)
		}
		return users, nil
	}
	return []int{}, nil
}
