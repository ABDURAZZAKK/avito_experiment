package service

import (
	"context"
	"fmt"

	"github.com/ABDURAZZAKK/avito_experiment/internal/entity"
	"github.com/ABDURAZZAKK/avito_experiment/internal/repo"
	"github.com/ABDURAZZAKK/avito_experiment/internal/repo/repoerrs"
)

type SegmentService struct {
	segmentRepo       repo.Segment
	usersSegmentsRepo repo.UsersSegments
}

func NewSegmentService(segmentRepo repo.Segment, usersSegmentsRepo repo.UsersSegments) *SegmentService {
	return &SegmentService{
		segmentRepo:       segmentRepo,
		usersSegmentsRepo: usersSegmentsRepo,
	}
}

func (s *SegmentService) GetBySlug(ctx context.Context, slug string) (entity.Segment, error) {
	segment, err := s.segmentRepo.GetBySlug(ctx, slug)
	if err != nil {
		if err == repoerrs.ErrNotFound {
			return entity.Segment{}, ErrNotFound
		}
		return entity.Segment{}, fmt.Errorf("SegmentService.GetBySlug - segmentRepo.GetBySlug: %v", err)
	}
	return segment, nil
}

func (s *SegmentService) Create(ctx context.Context, slug string) (string, error) {
	slug, err := s.segmentRepo.Create(ctx, slug)
	if err != nil {
		if err == repoerrs.ErrAlreadyExists {
			return "", ErrUserAlreadyExists
		}
		return "", fmt.Errorf("SegmentService.Create - segmentRepo.Create: %v", err)
	}
	return slug, nil
}
func (s *SegmentService) CreateAll(ctx context.Context, slugs []string) error {
	err := s.segmentRepo.CreateAll(ctx, slugs)
	if err != nil {
		if err == repoerrs.ErrAlreadyExists {
			return ErrAlreadyExists
		}
		return fmt.Errorf("SegmentService.CreateAll - segmentRepo.CreateAll: %v", err)
	}
	return nil
}

func (s *SegmentService) Delete(ctx context.Context, slug string) (string, error) {
	_, err := s.segmentRepo.GetBySlug(ctx, slug)
	if err != nil {
		if err == repoerrs.ErrNotFound {
			return "", ErrNotFound
		}
		return "", fmt.Errorf("SegmentService.Delete - segmentRepo.GetBySlug: %v", err)
	}
	s_slug, err := s.segmentRepo.Delete(ctx, slug)
	if err != nil {
		return "", fmt.Errorf("SegmentService.Delete - segmentRepo.Delete: %v", err)
	}

	return s_slug, nil
}
