CREATE TABLE IF NOT EXISTS collections (
    id           UUID NOT NULL PRIMARY KEY UNIQUE DEFAULT (UUID_V7()),
    public_id    UUID NOT NULL UNIQUE DEFAULT (UUID_V4()),
    workspace_id UUID NOT NULL,
    name         VARCHAR(255) NOT NULL,
    description  TEXT,
    category     VARCHAR(50),
    is_public    BOOLEAN NOT NULL DEFAULT FALSE,
    color        VARCHAR(7),
    icon         VARCHAR(100),
    public_stats BOOLEAN DEFAULT FALSE,
    pinned       BOOLEAN DEFAULT FALSE,
    deleted_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
    ON UPDATE CURRENT_TIMESTAMP,

    CONSTRAINT fk_collections_workspace_id
    FOREIGN KEY (workspace_id)
    REFERENCES workspaces (id)
    ON DELETE CASCADE
);

CREATE INDEX idx_collections_workspace_id_deleted_at
ON collections (workspace_id, deleted_at);
