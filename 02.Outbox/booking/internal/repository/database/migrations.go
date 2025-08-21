package database

import (
	"fmt"
	"strings"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"

	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// осуществление миграций
func (db *DataBase) Migrate() error {
	// получение драйвера
	driver, err := postgres.WithInstance(db.connection, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("failed to create migrate driver: %w", err)
	}

	// инициализация миграций
	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"postgres", driver)
	if err != nil {
		return fmt.Errorf("failed to init migrations: %w", err)
	}

	// Просто применяем миграции без сложной логики
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		// Если dirty state, используем force для текущей версии
		if strings.Contains(err.Error(), "Dirty database version") {
			version, _, verErr := m.Version()
			if verErr != nil {
				return fmt.Errorf("failed to get version: %w", verErr)
			}
			if forceErr := m.Force(int(version)); forceErr != nil {
				return fmt.Errorf("failed to force version %d: %w", version, forceErr)
			}
			// Повторяем миграцию
			return m.Up()
		}
		return fmt.Errorf("failed to apply migrations: %w", err)
	}

	return nil

}

// откат миграций
func (d *DataBase) RollbackLastMigration() error {
	// получение драйвера
	driver, err := postgres.WithInstance(d.connection, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("failed to create migrate driver: %w", err)
	}

	// инициализация миграций
	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"postgres", driver)
	if err != nil {
		return fmt.Errorf("failed to init migrations: %w", err)
	}

	// Откатываем на 1 шаг назад
	if err := m.Down(); err != nil {
		return fmt.Errorf("failed to rollback migration: %w", err)
	}

	return nil
}
