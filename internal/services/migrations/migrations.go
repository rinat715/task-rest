package migrations

import (
	"database/sql"
	r "go_rest/internal/repositories"
)

type MigrationService struct {
	db *sql.DB
}

func (m *MigrationService) Make(path string) error {
	return r.MakeMigrate(m.db, path)
}

func NewMigrationService(db *sql.DB) *MigrationService {
	return &MigrationService{
		db: db,
	}
}
