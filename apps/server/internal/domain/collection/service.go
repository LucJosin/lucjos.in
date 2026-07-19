package collection

import "context"

type Service interface {
	List(ctx context.Context) ([]Collection, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{
		repo: repo,
	}
}

func (s *service) List(ctx context.Context) ([]Collection, error) {
	return s.repo.List(ctx)
}
