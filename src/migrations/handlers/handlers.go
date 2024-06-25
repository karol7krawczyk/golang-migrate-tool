package handlers

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/Karol7Krawczyk/golang-migrate/migrations/config"
)

type Migration struct {
	Migration string
	AppliedAt time.Time
}

func HandleCommand(db *sql.DB, config config.Config) {
	switch {
	case config.Commands.History:
		handleHistoryCommand(db, config)
	case config.Commands.New:
		handleNewCommand(db, config)
	case config.Commands.Up:
		handleUpCommand(db, config)
	case config.Commands.Down:
		handleDownCommand(db, config)
	case config.Commands.Status != "":
		handleStatusCommand(db, config)
	case config.Commands.Create:
		handleCreateCommand(config)
	default:
		log.Fatalf("No valid command specified")
	}
}

func handleCreateCommand(config config.Config) {
	timestamp := time.Now().Format("20060102150405")
	migrationName := string(timestamp)
	description := strings.ReplaceAll(config.Commands.Desc, " ", "_")
	migrationDir := filepath.Join(config.Path, migrationName)

	if err := os.Mkdir(migrationDir, 0755); err != nil {
		log.Fatalf("Error creating migration directory: %v", err)
	}

	if err := os.WriteFile(filepath.Join(migrationDir, "up.sql"), []byte("-- Write your 'up' SQL here\n"), 0644); err != nil {
		log.Fatalf("Error creating 'up' migration file: %v", err)
	}

	if err := os.WriteFile(filepath.Join(migrationDir, "down.sql"), []byte("-- Write your 'down' SQL here\n"), 0644); err != nil {
		log.Fatalf("Error creating 'down' migration file: %v", err)
	}

	if config.Commands.Script {
		if err := os.WriteFile(filepath.Join(migrationDir, "up.sh"), []byte("echo 'Migration: "+migrationName+", bash script up'\n"), 0755); err != nil {
			log.Fatalf("Error creating script up.sh migration: %v", err)
		}

		if err := os.WriteFile(filepath.Join(migrationDir, "down.sh"), []byte("echo 'Migration: "+migrationName+", bash script down'\n"), 0755); err != nil {
			log.Fatalf("Error creating script down.sh migration: %v", err)
		}
	}

	if description != "" {
		if err := os.WriteFile(filepath.Join(migrationDir, description+".txt"), []byte("## "+config.Commands.Desc+"\n"), 0644); err != nil {
			log.Fatalf("Error creating description migration file: %v", err)
		}
	}

	fmt.Printf("Successfully created new migration: %s\n", migrationName)
}

func handleStatusCommand(db *sql.DB, config config.Config) {
	migration := config.Commands.Status
	historyMigrations, err := LoadHistoryMigrations(db, config)
	if err != nil {
		log.Fatalf("Error querying migrations: %v", err)
	}

	for _, m := range historyMigrations {
		if m.Migration == migration {
			fmt.Printf("Migration: %s, Applied At: %s\n", m.Migration, m.AppliedAt)
			return
		}
	}

	fmt.Printf("Migration %s not found in history\n", migration)
}

func handleHistoryCommand(db *sql.DB, config config.Config) {
	historyMigrations, err := LoadHistoryMigrations(db, config)
	if err != nil {
		log.Fatalf("Error querying migrations: %v", err)
	}

	fmt.Println("History of Migrations:")
	for _, m := range historyMigrations {
		fmt.Printf("Migration: %s, Applied At: %s\n", m.Migration, m.AppliedAt)
	}

	if len(historyMigrations) == 0 {
		fmt.Println("Migrations not added yet! You can check if there are new migrations")
		handleNewCommand(db, config)
	}
}

func handleNewCommand(db *sql.DB, config config.Config) {
	fmt.Println("New migrations to add:")
	newMigrations := loadNewMigrations(db, config)

	for _, migration := range newMigrations {
		fmt.Printf("Migration: %s, will be added\n", migration)
	}

	if len(newMigrations) == 0 {
		fmt.Println("There is nothing to add!")
	}
}

func handleUpCommand(db *sql.DB, config config.Config) {
	fmt.Println("Migrations to add:")
	newMigrations := loadNewMigrations(db, config)

	for _, migration := range newMigrations {
		if config.Commands.Steps > 0 || config.Commands.Steps < 0 {
			content := loadContent(filepath.Join(config.Path, migration, "up.sql"))

			scriptPath := filepath.Join(config.Path, migration, "up.sh")
			_, err := os.Stat(scriptPath)
			if !os.IsNotExist(err) {
				scriptPath := filepath.Join(config.Path, migration, "up.sh")
				if err := runBashScript(scriptPath, config); err != nil {
					log.Fatalf("Failed to run pre-migration script: %v", err)
				}
			}

			err = RunQueriesInTransaction(db, SplitSQLQueries(content))
			if err != nil {
				log.Fatalf("Error applying migration: %v %v", migration, err)
			}

			if config.Commands.Debug {
				fmt.Printf("-- DEBUG SQL: %s", content)
			}

			err = AddMigration(db, config, migration)
			if err != nil {
				log.Fatalf("Error adding migration:%v %v", migration, err)
				println(content)
			}

			fmt.Printf("Successfully applied migration: %s\n", migration)
			config.Commands.Steps = config.Commands.Steps - 1
		}
	}

	if len(newMigrations) == 0 {
		fmt.Println("There are no new migrations to apply.")
	}
}

func handleDownCommand(db *sql.DB, config config.Config) {
	fmt.Println("Migrations to remove:")
	historyMigrations, err := LoadHistoryMigrations(db, config)
	if err != nil {
		log.Fatalf("Error querying migrations: %v", err)
	}

	if config.Commands.Steps > 0 {
		historyMigrations = getMigrationsWithSteps(config.Commands.Steps, historyMigrations)
	}

	for _, m := range historyMigrations {
		content := loadContent(filepath.Join(config.Path, m.Migration, "down.sql"))

		err := RunQueriesInTransaction(db, SplitSQLQueries(content))
		if err != nil {
			log.Fatalf("Error run sql migration: %v", err)
		}

		if config.Commands.Debug {
			fmt.Printf("-- DEBUG SQL: %s", content)
		}

		err = RemoveMigration(db, config, m.Migration)
		if err != nil {
			log.Fatalf("Error remove migration: %v", err)
		}

		scriptPath := filepath.Join(config.Path, m.Migration, "down.sh")
		_, err = os.Stat(scriptPath)
		if !os.IsNotExist(err) {
			scriptPath := filepath.Join(config.Path, m.Migration, "down.sh")
			if err := runBashScript(scriptPath, config); err != nil {
				log.Fatalf("Failed to run script: %v", err)
			}
		}

		fmt.Printf("Migration '%s' has been successfully removed.\n", m.Migration)
	}

	if len(historyMigrations) == 0 {
		fmt.Println("There is nothing to remove!")
	}
}

func getMigrationsWithSteps(steps int, migrations []Migration) []Migration {
	n := len(migrations)
	if n <= steps {
		return migrations
	}
	return migrations[n-steps:]
}

func SplitSQLQueries(sqlQueries string) []string {
	queries := strings.Split(sqlQueries, ";")
	var cleanedQueries []string

	for _, query := range queries {
		cleanQuery := strings.TrimSpace(query)
		if cleanQuery != "" {
			cleanedQueries = append(cleanedQueries, cleanQuery)
		}
	}

	return cleanedQueries
}

func loadContent(filePath string) string {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		log.Fatalf("File %s does not exist", filePath)
	}

	content, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatalf("Error reading file %s: %v", filePath, err)
	}

	return string(content)
}

func RunQueriesInTransaction(db *sql.DB, queries []string) error {
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("error beginning transaction: %v", err)
	}
	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				log.Printf("Error rolling back transaction: %v", rbErr)
			}
			return
		}
		if commitErr := tx.Commit(); commitErr != nil {
			log.Printf("Error committing transaction: %v", commitErr)
		}
	}()

	for _, query := range queries {
		_, err := tx.Exec(query)
		if err != nil {
			return fmt.Errorf("error executing query: %v", err)
		}
	}

	return nil
}

func AddMigration(db *sql.DB, config config.Config, migration string) error {
	query := fmt.Sprintf("INSERT INTO %s (migration, applied_at) VALUES (?, ?)", config.TableName)

	record := Migration{
		Migration: migration,
		AppliedAt: time.Now(),
	}

	switch config.DBType {
	case "mysql":
		_, err := db.Exec(query, record.Migration, record.AppliedAt)
		if err != nil {
			return fmt.Errorf("error executing query: %v", err)
		}
	case "sqlite":
		_, err := db.Exec(query, record.Migration, record.AppliedAt.Format("2006-01-02 15:04:05"))
		if err != nil {
			return fmt.Errorf("error executing query: %v", err)
		}
	case "postgres":
		_, err := db.Exec(query, record.Migration, record.AppliedAt)
		if err != nil {
			return fmt.Errorf("error executing query: %v", err)
		}
	default:
		return fmt.Errorf("unsupported database type: %s", config.DBType)
	}

	return nil
}

func RemoveMigration(db *sql.DB, config config.Config, migration string) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE migration = ?", config.TableName)

	_, err := db.Exec(query, migration)
	if err != nil {
		return fmt.Errorf("error executing query: %v", err)
	}

	return nil
}

func LoadHistoryMigrations(db *sql.DB, config config.Config) ([]Migration, error) {
	query := fmt.Sprintf("SELECT migration, applied_at FROM %s ORDER BY migration ASC", config.TableName)

	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error querying migrations: %v", err)
	}
	defer rows.Close()

	var migrations []Migration
	for rows.Next() {
		var m Migration
		var appliedAtStr string
		if err := rows.Scan(&m.Migration, &appliedAtStr); err != nil {
			return nil, fmt.Errorf("error scanning migration row: %v", err)
		}

		switch config.DBType {
		case "mysql", "postgres":
			m.AppliedAt, err = time.Parse("2006-01-02 15:04:05", appliedAtStr)
		case "sqlite":
			m.AppliedAt, err = time.Parse(time.RFC3339, appliedAtStr)
		default:
			return nil, fmt.Errorf("unsupported database type: %s", config.DBType)
		}

		if err != nil {
			return nil, fmt.Errorf("error parsing applied_at timestamp: %v", err)
		}

		migrations = append(migrations, m)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over migration rows: %v", err)
	}

	return migrations, nil
}

func loadNewMigrations(db *sql.DB, config config.Config) []string {
	entries, err := os.ReadDir(config.Path)
	if err != nil {
		log.Fatalf("Error reading migration directory: %v", err)
	}

	historyMigrations, err := LoadHistoryMigrations(db, config)
	if err != nil {
		log.Fatalf("Error loading migration history: %v", err)
	}

	historySet := make(map[string]struct{})
	for _, m := range historyMigrations {
		historySet[m.Migration] = struct{}{}
	}

	var newMigrations []string
	for _, entry := range entries {
		if entry.IsDir() {
			_, exists := historySet[entry.Name()]
			if !exists {
				newMigrations = append(newMigrations, entry.Name())
			}
		}
	}

	sort.Strings(newMigrations)
	return newMigrations
}

func runBashScript(scriptPath string, config config.Config) error {
	cmd := exec.Command("/bin/bash", scriptPath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("error running script: %v\nOutput: %s", err, string(output))
	}

	if config.Commands.Debug {
		fmt.Printf("-- DEBUG SCRIPT: %s", string(output))
	}

	return nil
}
