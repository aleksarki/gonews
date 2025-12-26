package bootstrap

import (
	"fmt"
	"gonews/save_service/config"
	"gonews/save_service/internal/storage/pgstorage"
	"log"
)

func InitPGStorage(cfg *config.Config) *pgstorage.PGStorage {
	connectionString := fmt.Sprintf(
		"postgres://%s:%s@%s:%d:%s",
		cfg.Database.Username,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.DBName,
	)
	storage, err := pgstorage.NewPgstorage(connectionString)
	if err != nil {
		log.Panic(fmt.Sprintf("db initialization error, %v", err))
		panic(err)
	}

	return storage
}
