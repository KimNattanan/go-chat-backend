package app

import (
	"log"

	authEntity "github.com/KimNattanan/go-chat-backend/internal/auth/entity"
	"github.com/KimNattanan/go-chat-backend/internal/platform/config"
	profileEntity "github.com/KimNattanan/go-chat-backend/internal/profile/entity"
	"github.com/KimNattanan/go-chat-backend/pkg/postgres"
	"gorm.io/gorm"
)

func init() {
	cfg, err := config.NewConfig("")
	if err != nil {
		return
	}
	pg, err := postgres.New(cfg.DB.DSN)
	if err != nil {
		return
	}
	defer pg.Close()

	autoMigrateAll(pg.DB, cfg.App.ENV)
}

func autoMigrateAll(db *gorm.DB, env string) {
	if env == "test" {
		if err := db.Migrator().DropTable(
			&authEntity.User{},
			&profileEntity.Profile{},
		); err != nil {
			log.Fatalf("Migrate: drop table error: %s", err)
		}
	}
	if err := db.Migrator().AutoMigrate(
		&authEntity.User{},
		&profileEntity.Profile{},
	); err != nil {
		log.Fatalf("Migrate: auto migrate failed: %s", err)
	}
}
