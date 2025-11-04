package main

import (
	"log/slog"
	"sysken-pay-api/app/config"
	"sysken-pay-api/app/server"
)

func main() {
	db, err := server.NewDB(config.MySQLConfig())
	if err != nil {
		slog.Error("failed to create db", "err", err)
	}
	defer db.Close()

	if err = server.Run(db); err != nil {
		slog.Error("failed to run server", "err", err)
	}
}
