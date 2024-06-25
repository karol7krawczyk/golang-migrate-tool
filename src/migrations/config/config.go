package config

import (
	"flag"
	"os"
)

type Config struct {
	User      string
	Passwd    string
	TableName string
	Addr      string
	Port      string
	DBName    string
	Path      string
	DBType    string
	Commands  Commands
}

type Commands struct {
	History bool
	New     bool
	Up      bool
	Down    bool
	Debug   bool
	Step    bool
	Steps   int
	Status  string
	Create  bool
	Script  bool
	Desc    string
}

func ParseFlags() Config {
	var config Config

	flag.StringVar(&config.User, "db-user", os.Getenv("DB_USER"), "Database user")
	flag.StringVar(&config.Passwd, "db-password", os.Getenv("DB_PASSWORD"), "Database password")
	flag.StringVar(&config.Addr, "db-host", os.Getenv("DB_HOST"), "Database host address")
	flag.StringVar(&config.Port, "db-port", os.Getenv("DB_PORT"), "Database port")
	flag.StringVar(&config.DBName, "db-name", os.Getenv("DB_NAME"), "Database name")
	flag.StringVar(&config.TableName, "db-table", os.Getenv("DB_TABLE"), "Migration table")
	flag.StringVar(&config.Path, "path", os.Getenv("MIGRATION_PATH"), "Migration dir")
	flag.StringVar(&config.DBType, "db-type", os.Getenv("DB_TYPE"), "Database type (mysql, sqlite, postgres)")

	flag.StringVar(&config.Commands.Status, "status", "", "Check the status of a specific migration")
	flag.StringVar(&config.Commands.Desc, "desc", "", "Create a description of empty migration")
	flag.BoolVar(&config.Commands.Script, "script", false, "Create a bash scripts of empty migration")
	flag.BoolVar(&config.Commands.Create, "create", false, "Create a files of empty migration")
	flag.BoolVar(&config.Commands.History, "history", false, "Display migration history")
	flag.BoolVar(&config.Commands.New, "new", false, "Display upcoming migrations")
	flag.BoolVar(&config.Commands.Up, "up", false, "Apply new migrations")
	flag.BoolVar(&config.Commands.Down, "down", false, "Revert migrations")
	flag.BoolVar(&config.Commands.Debug, "debug", false, "Debug migrations")
	flag.BoolVar(&config.Commands.Step, "step", false, "Only one step migration")
	flag.IntVar(&config.Commands.Steps, "steps", -1, "Number of steps")

	flag.Parse()

	if config.Commands.Step {
		config.Commands.Steps = 1
		println(config.Commands.Steps)
	}

	return config
}
