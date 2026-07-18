CREATE TABLE IF NOT EXISTS workspace_users (
    workspace_id UUID NOT NULL,
    user_id      UUID NOT NULL,
    invited_by   UUID,
    role         VARCHAR(10) NOT NULL DEFAULT 'viewer' CHECK (
        role IN ('admin', 'member', 'viewer')
    ),
    deleted_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
    ON UPDATE CURRENT_TIMESTAMP,

    PRIMARY KEY (workspace_id, user_id),

    CONSTRAINT fk_workspace_users_workspace_id
    FOREIGN KEY (workspace_id)
    REFERENCES workspaces (id)
    ON DELETE CASCADE,

    CONSTRAINT fk_workspace_users_user_id
    FOREIGN KEY (user_id)
    REFERENCES users (id)
    ON DELETE CASCADE
);
