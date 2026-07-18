CREATE TABLE IF NOT EXISTS workspaces (
    id          UUID NOT NULL PRIMARY KEY UNIQUE DEFAULT (UUID_V7()),
    public_id   UUID NOT NULL UNIQUE DEFAULT (UUID_V4()),
    domain_id   UUID NOT NULL,
    slug        VARCHAR(255) NOT NULL UNIQUE,
    name        VARCHAR(255) NOT NULL,
    description TEXT,
    color       VARCHAR(7) NOT NULL DEFAULT '#000000',
    deleted_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
    ON UPDATE CURRENT_TIMESTAMP,

    CONSTRAINT fk_domains_workspace_id
    FOREIGN KEY (domain_id)
    REFERENCES domains (id)
    ON DELETE CASCADE
);
