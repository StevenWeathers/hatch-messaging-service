# Developer Guide

This document outlines the development setup and workflow for this project. We use Go 1.24 with the tools feature along with database-related tools.

## Prerequisites

- [Go 1.24](https://go.dev/dl/) or later
- [postgres 16](https://www.postgresql.org/) or later
- [Github](https://github.com) Account to fork and create PRs with

## Project Setup

1. Clone the repository:
   ```bash
   git clone https://github.com/StevenWeathers/hatch-messaging-service.git
   cd hatch-messaging-service
   ```

2. Install project dependencies:
   ```bash
   go mod download
   ```

## Database Management

### Database Migrations with Goose

We use [Goose](https://github.com/pressly/goose) to manage database schema migrations.

#### Creating a New Migration

```bash
go tool github.com/pressly/goose/v3/cmd/goose create -dir internal/db/migrations add_users_table sql
```

This creates a new migration file in the `internal/db/migrations` directory with a timestamp prefix.

#### Writing Migrations

Edit the newly created migration file. Each migration consists of an `Up` function (for applying the change) and a `Down` function (for reverting it).

Example:

```sql
-- +goose Up
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(255) NOT NULL UNIQUE,
    email VARCHAR(255) NOT NULL UNIQUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- +goose Down
DROP TABLE users;
```

#### Applying Migrations

Migrations are auto applied upon application start up, just build and run the application.