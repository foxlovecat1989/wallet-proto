# Database Migrations

This directory contains database migration files for the user service. The migrations are managed using [golang-migrate](https://github.com/golang-migrate/migrate).

## Migration Files

Migration files follow the naming convention: `{version}_{description}.{up|down}.sql`

- `{version}`: Sequential version number (e.g., 000001, 000002)
- `{description}`: Human-readable description of the migration
- `{up|down}`: Direction of the migration (up for applying, down for rolling back)

## Current Migrations

- `000001_init_schema.up.sql` / `000001_init_schema.down.sql`: Initial database schema setup

## Usage

### Using Makefile Commands

```bash
# Run all pending migrations
DATABASE_URL="postgres://user:pass@localhost:5432/dbname?sslmode=disable" make migrate-up

# Rollback the last migration
DATABASE_URL="postgres://user:pass@localhost:5432/dbname?sslmode=disable" make migrate-down

# Rollback multiple migrations
DATABASE_URL="postgres://user:pass@localhost:5432/dbname?sslmode=disable" STEPS=3 make migrate-down

# Check migration status
DATABASE_URL="postgres://user:pass@localhost:5432/dbname?sslmode=disable" make migrate-status

# Create a new migration
NAME="add_user_profile" make migrate-create
```

### Using the CLI Tool Directly

```bash
# Build the migration tool
make build-migrate

# Run migrations
./bin/migrate -database "postgres://user:pass@localhost:5432/dbname?sslmode=disable" -action up

# Rollback migrations
./bin/migrate -database "postgres://user:pass@localhost:5432/dbname?sslmode=disable" -action down -steps 1

# Check status
./bin/migrate -database "postgres://user:pass@localhost:5432/dbname?sslmode=disable" -action status
```

### Automatic Migrations

The application automatically runs migrations on startup. This ensures the database schema is always up to date.

## Best Practices

1. **Always create both up and down migrations**: Every migration should have a corresponding rollback.
2. **Use descriptive names**: Migration names should clearly describe what the migration does.
3. **Test migrations**: Always test both up and down migrations before committing.
4. **Keep migrations small**: Each migration should make a single logical change.
5. **Use transactions**: Wrap migration logic in transactions when possible.
6. **Version control**: Always commit migration files to version control.

## Creating New Migrations

1. Use the make command to create migration files:
   ```bash
   NAME="add_user_profile" make migrate-create
   ```

2. Edit the generated `.up.sql` and `.down.sql` files with your schema changes.

3. Test the migration:
   ```bash
   # Test up migration
   DATABASE_URL="..." make migrate-up
   
   # Test down migration
   DATABASE_URL="..." make migrate-down
   ```

4. Commit the migration files to version control.

## Migration Schema

The migration system creates a `schema_migrations` table in your database to track which migrations have been applied. This table is automatically managed by golang-migrate.

## Troubleshooting

### Migration Already Applied
If you get an error about a migration already being applied, you can force the migration:
```bash
./bin/migrate -database "..." -action force -version 1
```

### Dirty Database
If the database is in a "dirty" state (migration failed partway through), you can fix it:
```bash
./bin/migrate -database "..." -action force -version 1
```

### Connection Issues
Make sure your DATABASE_URL is correct and the database is accessible. The URL format is:
```
postgres://username:password@host:port/database?sslmode=disable
```
