package main

import (
	"database/sql"
	"fmt"
	"log/slog"
	"os"

	_ "github.com/go-sql-driver/mysql"
	migrate "github.com/rubenv/sql-migrate"
	"github.com/zuxt268/berry/internal/config"
)

func main() {
	direction := "up"
	if len(os.Args) > 1 {
		direction = os.Args[1]
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&charset=utf8mb4&loc=Asia%%2FTokyo",
		config.Env.DBUser, config.Env.DBPassword,
		config.Env.DBHost, config.Env.DBPort, config.Env.DBName)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		slog.Error("failed to open database", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	migrations := &migrate.FileMigrationSource{Dir: "db/migrations"}

	var n int
	switch direction {
	case "up":
		n, err = migrate.Exec(db, "mysql", migrations, migrate.Up)
	case "down":
		n, err = migrate.Exec(db, "mysql", migrations, migrate.Down)
	default:
		slog.Error("unknown direction, use 'up' or 'down'", "direction", direction)
		os.Exit(1)
	}

	if err != nil {
		slog.Error("migration failed", "direction", direction, "error", err)
		os.Exit(1)
	}

	slog.Info("migration completed", "direction", direction, "applied", n)
}
