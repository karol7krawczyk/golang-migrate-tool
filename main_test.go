package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/Karol7Krawczyk/golang-migrate/migrations/config"
	"github.com/Karol7Krawczyk/golang-migrate/migrations/db"
	"github.com/Karol7Krawczyk/golang-migrate/migrations/handlers"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

var testConfig config.Config

func setupTestDB(t *testing.T) *sql.DB {
	testConfig = config.Config{
		User:      os.Getenv("DB_USER"),
		Passwd:    os.Getenv("DB_PASSWORD"),
		TableName: os.Getenv("DB_TABLE"),
		Addr:   os.Getenv("DB_HOST") + ":" + os.Getenv("DB_HOST"), // Change this based on your database configuration
		DBName: os.Getenv("DB_NAME"),
		Path:   os.Getenv("MIGRATION_PATH"),
		DBType: os.Getenv("DB_TYPE"), // Change this to mysql or postgres as needed
	}

	database, err := db.GetConnection(testConfig)
	if err != nil {
		t.Fatalf("Error connecting to the database: %v", err)
	}

	err = db.PrepareMigrationTable(database, testConfig)
	if err != nil {
		t.Fatalf("Error preparing migration table: %v", err)
	}

	return database
}

func teardownTestDB(db *sql.DB) {
	query := fmt.Sprintf("DROP TABLE IF EXISTS %s", testConfig.TableName)
	_, err := db.Exec(query)
	if err != nil {
		log.Printf("Error dropping test table: %v", err)
	}

	db.Close()
}

func TestMain(m *testing.M) {
	// Run tests
	code := m.Run()
	os.Exit(code)
}

func TestGetConnection(t *testing.T) {
	db := setupTestDB(t)
	defer teardownTestDB(db)

	if err := db.Ping(); err != nil {
		t.Fatalf("Failed to ping database: %v", err)
	}
}

func TestPrepareMigrationTable(t *testing.T) {
	db := setupTestDB(t)
	defer teardownTestDB(db)

	query := fmt.Sprintf("SELECT 1 FROM %s LIMIT 1", testConfig.TableName)
	_, err := db.Exec(query)
	if err != nil {
		t.Fatalf("Migration table not prepared correctly: %v", err)
	}
}

func TestAddAndRemoveMigration(t *testing.T) {
	db := setupTestDB(t)
	defer teardownTestDB(db)

	migration := "22220101120001"
	err := handlers.AddMigration(db, testConfig, migration)
	if err != nil {
		t.Fatalf("Failed to add migration: %v", err)
	}

	historyMigrations, err := handlers.LoadHistoryMigrations(db, testConfig)
	if err != nil {
		t.Fatalf("Failed to load history migrations: %v", err)
	}

	found := false
	for _, m := range historyMigrations {
		if m.Migration == migration {
			found = true
			break
		}
	}

	if !found {
		t.Fatalf("Migration %s not found in history", migration)
	}

	err = handlers.RemoveMigration(db, testConfig, migration)
	if err != nil {
		t.Fatalf("Failed to remove migration: %v", err)
	}

	historyMigrations, err = handlers.LoadHistoryMigrations(db, testConfig)
	if err != nil {
		t.Fatalf("Failed to load history migrations: %v", err)
	}

	found = false
	for _, m := range historyMigrations {
		if m.Migration == migration {
			found = true
			break
		}
	}

	if found {
		t.Fatalf("Migration %s still found in history after removal", migration)
	}
}

func TestRunQueriesInTransaction(t *testing.T) {
	db := setupTestDB(t)
	defer teardownTestDB(db)

	queries := []string{
		fmt.Sprintf("INSERT INTO %s (migration, applied_at) VALUES ('20240208100231', '%s')", testConfig.TableName, time.Now().Format("2006-01-02 15:04:05")),
		fmt.Sprintf("INSERT INTO %s (migration, applied_at) VALUES ('20240209110543', '%s')", testConfig.TableName, time.Now().Format("2006-01-02 15:04:05")),
	}

	err := handlers.RunQueriesInTransaction(db, queries)
	if err != nil {
		t.Fatalf("Failed to run queries in transaction: %v", err)
	}

	historyMigrations, err := handlers.LoadHistoryMigrations(db, testConfig)
	if err != nil {
		t.Fatalf("Failed to load history migrations: %v", err)
	}

	if len(historyMigrations) != 2 {
		t.Fatalf("Expected 2 migrations, found %d", len(historyMigrations))
	}
}

func TestSplitSQLQueries(t *testing.T) {
	sqlQuery := "CREATE TABLE test1 (id INT); CREATE TABLE test2 (id INT);"
	expected := []string{
		"CREATE TABLE test1 (id INT)",
		"CREATE TABLE test2 (id INT)",
	}

	result := handlers.SplitSQLQueries(sqlQuery)
	if len(result) != len(expected) {
		t.Fatalf("Expected %d queries, got %d", len(expected), len(result))
	}

	for i, query := range result {
		if query != expected[i] {
			t.Fatalf("Expected query %s, got %s", expected[i], query)
		}
	}
}
