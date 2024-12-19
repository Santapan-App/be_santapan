package commands

import "database/sql"

type PostgresMenuCommandRepository struct {
	Conn *sql.DB
}

func NewPostgresMenuCommandRepository(Conn *sql.DB) *PostgresMenuCommandRepository {
	return &PostgresMenuCommandRepository{Conn}
}
