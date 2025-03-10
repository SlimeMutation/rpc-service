package database

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"os"
	"path/filepath"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/SlimeMutation/rpc-service/common/retry"
	"github.com/SlimeMutation/rpc-service/config"
)

type DB struct {
	gorm *gorm.DB
	Keys KeysDB
}

func NewDB(ctx context.Context, dbConf config.DBConfig) (*DB, error) {
	dsn := fmt.Sprintf("host=%s dbname=%s sslmode=disable", dbConf.Host, dbConf.Name)
	if dbConf.Port != 0 {
		dsn += fmt.Sprintf(" port=%d", dbConf.Port)
	}
	if dbConf.User != "" {
		dsn += fmt.Sprintf(" user=%s", dbConf.User)
	}
	if dbConf.Password != "" {
		dsn += fmt.Sprintf(" password=%s", dbConf.Password)
	}

	gormConfig := gorm.Config{
		SkipDefaultTransaction: true,
		CreateBatchSize:        3_000,
	}

	retryStrategy := &retry.ExponentialStrategy{Min: 1000, Max: 20_000, MaxJitter: 250}
	gorm1, err := retry.Do[*gorm.DB](context.Background(), 10, retryStrategy, func() (*gorm.DB, error) {
		gorm2, err := gorm.Open(postgres.Open(dsn), &gormConfig)
		if err != nil {
			return nil, fmt.Errorf("failed to connect to database: %w", err)
		}
		return gorm2, nil
	})
	if err != nil {
		return nil, err
	}
	db := &DB{
		gorm: gorm1,
		Keys: NewKeysDB(gorm1),
	}
	return db, nil
}

func (db *DB) Transaction(fn func(db *DB) error) error {
	return db.gorm.Transaction(func(tx *gorm.DB) error {
		txDB := &DB{
			gorm: tx,
			Keys: NewKeysDB(tx),
		}
		return fn(txDB)
	})
}

func (db *DB) Close() error {
	sql, err := db.gorm.DB()
	if err != nil {
		return err
	}
	return sql.Close()
}

func (db *DB) ExecuteSQLMigration(migrationsFolder string) error {
	err := filepath.Walk(migrationsFolder, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("Failed to process migration file: %s", path))
		}
		if info.IsDir() {
			return nil
		}
		fileContent, readErr := os.ReadFile(path)
		if readErr != nil {
			return errors.Wrap(readErr, fmt.Sprintf("Error reading SQL file: %s", path))
		}

		execErr := db.gorm.Exec(string(fileContent)).Error
		if execErr != nil {
			return errors.Wrap(execErr, fmt.Sprintf("Error executing SQL script: %s", path))
		}
		return nil
	})
	return err
}
