package commands

import (
	"log"

	"spotistats/internal/config"
	"spotistats/internal/database"

	"github.com/spf13/cobra"
)

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Run database migrations",
	Run:   runMigrate,
}

func runMigrate(cmd *cobra.Command, args []string) {
	cfg := config.Load()

	if err := database.Connect(cfg.DatabaseURL()); err != nil {
		log.Fatal("failed to connect to db:", err)
	}
	defer database.Close()

	if err := database.Migrate(); err != nil {
		log.Fatal("migration failed:", err)
	}

	log.Println("migrations completed successfully")
}
