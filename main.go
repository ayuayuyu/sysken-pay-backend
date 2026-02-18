package main

import (
	"database/sql"
	"log/slog"
	"sysken-pay-api/app/config"
	"sysken-pay-api/app/server"
	"time"
)

func main() {
	var db *sql.DB
	var err error

	maxRetries := 30
	waitTime := 2 * time.Second

	for i := 0; i < maxRetries; i++ {
		db, err = server.NewDB(config.MySQLConfig())

		if err == nil {
			if errPing := db.Ping(); errPing == nil {
				slog.Info("Successfully connected to database")
				break
			} else {
				db.Close()
				err = errPing
			}
		}

		slog.Warn("waiting for database...", "attempt", i+1, "err", err)
		time.Sleep(waitTime)
	}

	if err != nil {
		slog.Error("failed to create db after retries", "err", err)
		return
	}
	defer db.Close()

	if err = server.Run(db); err != nil {
		slog.Error("failed to run server", "err", err)
	}
}
