package system

import (
	"context"

	"github.com/gofrs/uuid/v5"
)

type Service interface {
	Info(ctx context.Context) (System, error)
	Configure(ctx context.Context, ownerUserID uuid.UUID) (System, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{
		repo: repo,
	}
}

func (s *service) Info(ctx context.Context) (System, error) {
	return s.repo.FindOne(ctx)
}

func (s *service) Configure(ctx context.Context, ownerUserID uuid.UUID) (System, error) {
	return s.repo.Create(ctx, ownerUserID)
}
