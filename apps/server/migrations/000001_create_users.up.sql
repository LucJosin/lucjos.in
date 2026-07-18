CREATE TABLE IF NOT EXISTS users (
    id             UUID NOT NULL PRIMARY KEY UNIQUE DEFAULT (UUID_V7()),
    first_name     VARCHAR(50) NOT NULL,
    last_name      VARCHAR(50),
    username       VARCHAR(50) NOT NULL UNIQUE,
    email          VARCHAR(255) NOT NULL UNIQUE,
    email_verified BOOLEAN DEFAULT FALSE,
    password       VARCHAR(255),
    avatar_url     VARCHAR(255),
    deleted_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
    ON UPDATE CURRENT_TIMESTAMP
);

CREATE INDEX idx_users_username_deleted_at
ON users (username, deleted_at);
