package domain

import (
	"context"
	"errors"
	"fmt"

	"github.com/lucjosin/qorv.in/internal/errs"
)

type Service interface {
	FindByDomain(ctx context.Context, domain string) (Domain, error)
	Create(ctx context.Context, entity Domain) (Domain, error)
	FindOrCreateByDomain(ctx context.Context, domain Domain) (Domain, bool, error)
}

type service struct {
	repo Repository
}

func NewService(repository Repository) Service {
	return &service{repo: repository}
}

func (s *service) FindByDomain(ctx context.Context, domain string) (Domain, error) {
	return s.repo.FindByDomain(ctx, domain)
}

func (s *service) Create(ctx context.Context, entity Domain) (Domain, error) {
	return s.repo.Create(ctx, entity)
}

func (s *service) FindOrCreateByDomain(ctx context.Context, entity Domain) (Domain, bool, error) {
	existingEntity, err := s.repo.FindByDomain(ctx, entity.Domain)
	if err == nil {
		return existingEntity, false, nil
	}
	if !errors.Is(err, errs.ErrNotFound) {
		return Domain{}, false, fmt.Errorf("checking domain existence: %w", err)
	}

	createdEntity, err := s.repo.Create(ctx, entity)
	if err == nil {
		return createdEntity, true, nil
	}

	if errors.Is(err, errs.ErrConflict) {
		// query the DB one more time to retrieve the recent data
		recentEntity, err := s.repo.FindByDomain(ctx, createdEntity.Domain)
		if err != nil {
			return Domain{}, false, fmt.Errorf("retrieving domain after conflict: %w", err)
		}
		return recentEntity, false, nil
	}

	return Domain{}, false, fmt.Errorf("creating domain: %w", err)
}
