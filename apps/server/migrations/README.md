# Database Migrations

This directory contains the SQL migration files managed by `golang-migrate/migrate`. 

All migrations are structured in matching pairs: 
- An `.up.sql` file to apply the change
- An `.down.sql` file to completely revert it.

## Filename Pattern

```text
{version}_{title}.up.sql
{version}_{title}.down.sql
```

- **`{version}`**: Sequential 64-bit integer, left-padded with zeros (e.g., `000001`).
- **`{title}`**: Snake_case description of the change (e.g., `create_users_table`).

> [!NOTE]
> When using the `golang-migrate/migrate` cli, **you don't need to specify the sequential number**, provide only the title.

### Creating New Migrations

Always use the CLI to generate the files with the correct sequential sequence:

```bash
migrate create -ext sql -dir migrations -seq <migration_name>
```

or, you can also use the [Makefile](../Makefile) migration helper:

```bash
make migrate/new NAME=<migration_name>
```

## Migration Patterns & Examples

Use the following SQL patterns to ensure migrations execute safely and can rollback cleanly.

### Create Table

- **`000001_create_users_table.up.sql`**
    ```sql
    CREATE TABLE IF NOT EXISTS users (
        id BIGSERIAL PRIMARY KEY,
        email VARCHAR(255) UNIQUE NOT NULL,
        password_hash TEXT NOT NULL,
        created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
    );
    ```
- **`000001_create_users_table.down.sql`**
    ```sql
    DROP TABLE IF EXISTS users;
    ```

### Add Column

- **000002_add_bio_to_users.up.sql**
    ```sql
    ALTER TABLE users ADD COLUMN IF NOT EXISTS bio TEXT;
    ```
- **000002_add_bio_to_users.down.sql**
    ```sql
    ALTER TABLE users DROP COLUMN IF EXISTS bio;
    ```

### Alter Column (Data Type / Constraints)

- **000003_alter_users_bio_length.up.sql**
    ```sql
    ALTER TABLE users ALTER COLUMN bio TYPE VARCHAR(500);
    ```
- **000003_alter_users_bio_length.down.sql**
    ```sql
    ALTER TABLE users ALTER COLUMN bio TYPE TEXT;
    ```

### Create Index

- **000004_add_email_index_to_users.up.sql**
    ```sql
    CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_users_email ON users(email);
    ```
- **000004_add_email_index_to_users.down.sql**
    ```sql
    DROP INDEX IF EXISTS idx_users_email;
    ```

## Notes

1. **Idempotency:** Use `IF NOT EXISTS` and `IF EXISTS` where possible to prevent the migration tool from crashing if a state was partially applied.
2. **One Feature Per Pair:** Do not combine unrelated schema changes (e.g., creating a `posts` table and modifying the `users` table) in a single migration file. 
3. **Immutability:** Never modify a migration file that has already been merged into production. Create a new sequential migration to change or undo it.
