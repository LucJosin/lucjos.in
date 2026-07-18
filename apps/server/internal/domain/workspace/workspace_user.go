package workspace

import (
	"fmt"
	"slices"
	"time"

	"github.com/gofrs/uuid/v5"
)

// WorkspaceRole represents a user's permission level within a workspace.
type WorkspaceRole string

const (
	// OwnerRole has full control, including billing and deletion.
	OwnerRole WorkspaceRole = "owner"

	// AdminRole can manage settings and members.
	AdminRole WorkspaceRole = "admin"

	// MemberRole can create and edit resources.
	MemberRole WorkspaceRole = "member"

	// ViewerRole has read-only access.
	ViewerRole WorkspaceRole = "viewer"
)

var ValidWorkspaceRoles = []WorkspaceRole{
	OwnerRole,
	AdminRole,
	MemberRole,
	ViewerRole,
}

type WorkspaceUser struct {
	WorkspaceID uuid.UUID
	UserID      uuid.UUID
	InvitedBy   *uuid.UUID
	Role        WorkspaceRole
	DeletedAt   *time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (wu *WorkspaceUser) ValidateAndPrepare() error {
	wu.Prepare()
	return wu.Validate()
}

func (wu *WorkspaceUser) Validate() error {
	if !slices.Contains(ValidWorkspaceRoles, wu.Role) {
		return fmt.Errorf("invalid workspace role: %s", wu.Role)
	}
	return nil
}

func (wu *WorkspaceUser) Prepare() {
}
