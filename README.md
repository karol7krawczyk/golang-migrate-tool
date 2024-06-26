# Golang Migrate

![License](https://img.shields.io/badge/license-MIT-blue.svg)

## Table of Contents
- [Overview](#overview)
- [Features](#features)
- [Installation](#installation)
- [Usage](#usage)
- [Configuration](#configuration)
- [Command-Line Flags](#command-line-flags)
- [Testing](#testing)
- [License](#license)
- [Contact](#contact)

## Overview
The `golang-migrate` project is designed to provide an efficient and easy-to-use solution for database migrations in Golang. It supports multiple databases and includes features such as version control, rollback, and detailed logging. Project is for people who need to run a bash script before executing the SQL code or are looking for migrations that execute the processes in scripts

## Features
- Supports multiple databases (e.g., MySQL, PostgreSQL, SQLite)
- Efficiently using either SQL files and bash scripts.
- Version control for migrations
- Rollback functionality
- Detailed logging
- Easy configuration

## Installation
To get started, clone the repository and install the required dependencies.

### Clone the repository
```bash
git clone https://github.com/Karol7Krawczyk/golang-migrate.git
cd golang-migrate
```

### Run in docker
```bash
make build
make up
```

### The binary file is available after building the docker
```bash
./migration -up
```

## Usage
Here's how you can use the migration tool:

- List of all commands: go run main.go -help
- Running Migrations: go run main.go -up
- Rolling Back Migrations: go run main.go -down
- Create Migration files: go run main.go -create

## Configuration
Configuration can be done via environment variables or command-line flags. Here’s an example configuration for your project:

```bash
    environment:
      DB_NAME: your-database
      DB_USER: your-username
      DB_PASSWORD: your-password
      DB_HOST: localhost
      DB_PORT: 3306
      DB_TYPE: mysql
      DB_TABLE: migrations
      MIGRATION_PATH: migrations/data
```

## Command-Line Flags
The migration tool accepts various command-line flags for configuration and commands. You can also use environment variables for configuration.

### Database Configuration
- `-db-user`: Database user (default: `DB_USER` environment variable)
- `-db-password`: Database password (default: `DB_PASSWORD` environment variable)
- `-db-host`: Database host address (default: `DB_HOST` environment variable)
- `-db-port`: Database port (default: `DB_PORT` environment variable)
- `-db-name`: Database name (default: `DB_NAME` environment variable)
- `-db-table`: Migration table (default: `DB_TABLE` environment variable)
- `-path`: Migration directory (default: `MIGRATION_PATH` environment variable)
- `-db-type`: Database type (mysql, sqlite, postgres) (default: `DB_TYPE` environment variable)

### Commands
- `-status`: Check the status of a specific migration
- `-desc`: Create a description of an empty migration
- `-script`: Create bash scripts for an empty migration
- `-create`: Create files for an empty migration
- `-history`: Display migration history
- `-new`: Display upcoming migrations
- `-up`: Apply new migrations
- `-down`: Revert migrations
- `-debug`: Debug migrations
- `-step`: Only one-step migration
- `-steps`: Number of steps for migration

### Example Usage
```bash
go run . -db-user=root -db-password=secret -db-host=localhost -db-port=3306 -db-name=migrations -up
go run . -create -script -desc="Create User table"
go run . -create
go run . -debug -down
go run . -debug -down -step
go run . -debug -up -steps=5
go run . -history
```

## Testing
To run the tests, use the following command after build docker-compose:

```bash
make up
make test
```

## License
This project is licensed under the MIT License - see the LICENSE file for details.
