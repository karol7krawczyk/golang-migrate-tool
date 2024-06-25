package db

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/Karol7Krawczyk/golang-migrate/migrations/config"
	"github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

func GetConnection(config config.Config) (*sql.DB, error) {
	switch config.DBType {
	case "mysql":
		cfg := mysql.Config{
			User:                 config.User,
			Passwd:               config.Passwd,
			Net:                  "tcp",
			Addr:                 config.Addr,
			DBName:               config.DBName,
			AllowNativePasswords: true,
		}

		db, err := sql.Open("mysql", cfg.FormatDSN())
		if err != nil {
			return nil, fmt.Errorf("error connecting to the database: %v", err)
		}
		if err := db.Ping(); err != nil {
			return nil, fmt.Errorf("error pinging the database: %v", err)
		}
		return db, nil
	case "postgres":
		psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
			config.Addr, config.Port, config.User, config.Passwd, config.DBName, "disable")

		db, err := sql.Open("postgres", psqlInfo)
		if err != nil {
			return nil, fmt.Errorf("error connecting to the database: %v", err)
		}
		if err := db.Ping(); err != nil {
			return nil, fmt.Errorf("error pinging the database: %v", err)
		}
		return db, nil
	case "sqlite":
		//const dbPath = "./your_database.db" // Replace with the actual path to your database file
		db, err := sql.Open("sqlite3", config.DBName)
		if err != nil {
			panic(err)
		}
		if err := db.Ping(); err != nil {
			return nil, fmt.Errorf("error pinging the database: %v", err)
		}
		return db, nil
	default:
		return nil, fmt.Errorf("unsupported database type: %s", config.DBType)
	}
}

func CloseConnection(db *sql.DB) {
	if err := db.Close(); err != nil {
		log.Printf("Error closing the database: %v", err)
	}
}

func PrepareMigrationTable(db *sql.DB, config config.Config) error {
	var query string
	switch config.DBType {
	case "mysql":
		query = fmt.Sprintf(`
        SELECT EXISTS (
            SELECT 1
            FROM   information_schema.tables
            WHERE  table_schema = DATABASE()
            AND    table_name = '%s'
        );`, config.TableName)
	case "sqlite":
		query = fmt.Sprintf(`
        SELECT EXISTS (
            SELECT 1
            FROM   sqlite_master
            WHERE  type='table'
            AND    name='%s'
        );`, config.TableName)
	case "postgres":
		query = fmt.Sprintf(`
        SELECT EXISTS (
            SELECT 1
            FROM   pg_tables
            WHERE  schemaname = 'public'
            AND    tablename = '%s'
        );`, config.TableName)
	default:
		return fmt.Errorf("unsupported database type: %s", config.DBType)
	}

	var exists bool
	if err := db.QueryRow(query).Scan(&exists); err != nil {
		return fmt.Errorf("error executing query: %v", err)
	}

	if !exists {
		var createTableQuery string
		switch config.DBType {
		case "mysql":
			createTableQuery = fmt.Sprintf(`
            CREATE TABLE %s (
                migration VARCHAR(255) NOT NULL PRIMARY KEY,
                applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
            );`, config.TableName)
		case "sqlite":
			createTableQuery = fmt.Sprintf(`
            CREATE TABLE %s (
                migration TEXT NOT NULL PRIMARY KEY,
                applied_at DATETIME DEFAULT CURRENT_TIMESTAMP
            );`, config.TableName)
		case "postgres":
			createTableQuery = fmt.Sprintf(`
            CREATE TABLE %s (
                migration VARCHAR(255) NOT NULL PRIMARY KEY,
                applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
            );`, config.TableName)
		default:
			return fmt.Errorf("unsupported database type: %s", config.DBType)
		}

		if _, err := db.Exec(createTableQuery); err != nil {
			return fmt.Errorf("error creating table: %v", err)
		}

		fmt.Printf("Table %s created successfully\n", config.TableName)
	}

	return nil
}
