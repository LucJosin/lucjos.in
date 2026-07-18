package workspace

import (
	"context"
	"errors"
	"fmt"

	"github.com/gofrs/uuid/v5"
	"github.com/lucjosin/qorv.in/internal/domain/user"
	"github.com/lucjosin/qorv.in/internal/errs"
)

type Service interface {
	FindByDomainID(ctx context.Context, domainID uuid.UUID) (Workspace, error)
	FindOrCreateByDomainID(ctx context.Context, entity Workspace) (Workspace, bool, error)
	Create(ctx context.Context, entity Workspace) (Workspace, error)

	AddUser(ctx context.Context, entity WorkspaceUser) error
}

type service struct {
	repo        Repository
	userService user.Service
}

func NewService(repo Repository, userService user.Service) Service {
	return &service{
		repo:        repo,
		userService: userService,
	}
}

func (s *service) FindByDomainID(ctx context.Context, domainID uuid.UUID) (Workspace, error) {
	return s.repo.FindByDomainID(ctx, domainID)
}

func (s *service) FindOrCreateByDomainID(ctx context.Context, entity Workspace) (Workspace, bool, error) {
	existingEntity, err := s.repo.FindByDomainID(ctx, entity.DomainID)
	if err == nil {
		return existingEntity, false, nil
	}
	if !errors.Is(err, errs.ErrNotFound) {
		return Workspace{}, false, fmt.Errorf("checking workspace existence: %w", err)
	}

	createdEntity, err := s.repo.Create(ctx, entity)
	if err == nil {
		return createdEntity, true, nil
	}

	if errors.Is(err, errs.ErrConflict) {
		// query the DB one more time to retrieve the recent data
		recentEntity, err := s.repo.FindByDomainID(ctx, createdEntity.DomainID)
		if err != nil {
			return Workspace{}, false, fmt.Errorf("retrieving workspace after conflict: %w", err)
		}
		return recentEntity, false, nil
	}

	return Workspace{}, false, fmt.Errorf("creating workspace: %w", err)
}

func (s *service) Create(ctx context.Context, entity Workspace) (Workspace, error) {
	err := entity.ValidateAndPrepare()
	if err != nil {
		return Workspace{}, err
	}

	created, err := s.repo.Create(ctx, entity)
	if err != nil {
		return Workspace{}, err
	}

	return created, nil
}

func (s *service) AddUser(ctx context.Context, entity WorkspaceUser) error {
	exists, err := s.repo.ExistsByID(ctx, entity.WorkspaceID)
	if err != nil {
		return err
	}
	if !exists {
		return errs.ErrNotFound
	}

	exists, err = s.userService.ExistsByID(ctx, entity.UserID)
	if err != nil {
		return err
	}
	if !exists {
		return errs.ErrNotFound
	}

	err = s.repo.CreateWorkspaceUser(ctx, entity)
	if err != nil {
		return fmt.Errorf("adding user to workspace: %w", err)
	}

	return nil
}
