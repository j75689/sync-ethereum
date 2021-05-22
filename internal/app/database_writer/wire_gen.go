// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//+build !wireinject

package database_writer

import (
	"sync-ethereum/internal/config"
	"sync-ethereum/internal/delivery/database_writer"
	"sync-ethereum/internal/repository/gorm"
	"sync-ethereum/internal/service/storage"
	"sync-ethereum/internal/wireset"
)

// Injectors from wire.go:

func Initialize(configPath string) (Application, error) {
	configConfig, err := config.NewConfig(configPath)
	if err != nil {
		return Application{}, err
	}
	logger, err := wireset.InitLogger(configConfig)
	if err != nil {
		return Application{}, err
	}
	mq, err := wireset.InitMQ(configConfig, logger)
	if err != nil {
		return Application{}, err
	}
	db, err := wireset.InitDatabase(configConfig, logger)
	if err != nil {
		return Application{}, err
	}
	storageRepository := gorm.NewStorageRepository(db)
	storageService := storage.NewStorageService(storageRepository)
	databaseWriter := database_writer.NewDatabaseWriter(configConfig, logger, mq, storageService)
	application := newApplication(logger, databaseWriter)
	return application, nil
}
