package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/gofrs/uuid/v5"
	"github.com/lucjosin/qorv.in/internal/errs"
)

type Service interface {
	FindByID(ctx context.Context, id uuid.UUID) (User, error)
	FindOrCreateByUsername(ctx context.Context, entity User) (User, bool, error)
	ExistsByID(ctx context.Context, id uuid.UUID) (bool, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{
		repo: repo,
	}
}

func (s *service) FindByID(ctx context.Context, id uuid.UUID) (User, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *service) FindOrCreateByUsername(ctx context.Context, entity User) (User, bool, error) {
	existingUser, err := s.repo.FindByUsername(ctx, entity.Username)
	if err == nil {
		return existingUser, false, nil
	}
	if !errors.Is(err, errs.ErrNotFound) {
		return User{}, false, fmt.Errorf("checking user existence: %w", err)
	}

	user, err := s.repo.Create(ctx, entity)
	if err == nil {
		return user, true, nil
	}

	if errors.Is(err, errs.ErrConflict) {
		// query the DB one more time to retrieve the recent data
		recentUser, err := s.repo.FindByUsername(ctx, user.Username)
		if err != nil {
			return User{}, false, fmt.Errorf("retrieving user after conflict: %w", err)
		}
		return recentUser, false, nil
	}

	return User{}, false, fmt.Errorf("creating user: %w", err)
}

func (s *service) ExistsByID(ctx context.Context, id uuid.UUID) (bool, error) {
	return s.repo.ExistsByID(ctx, id)
}
