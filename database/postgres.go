package database

import (
	"context"
	"database/sql"
	_ "github.com/lib/pq"
	"rest-ws/models"
)

type PostgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(url string) (*PostgresRepository, error) {
	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, err
	}
	return &PostgresRepository{db: db}, nil
}

func (p *PostgresRepository) InsertUser(ctx context.Context, user *models.User) error {
	_, err := p.db.ExecContext(ctx, "INSERT INTO users (email, password) VALUES ($1, $2)", user.Email, user.Password)
	return err
}

func (p *PostgresRepository) GetUserById(ctx context.Context, id int64) (*models.User, error) {
	row := p.db.QueryRowContext(ctx, "SELECT id, email FROM users WHERE id = $1", id)
	var user models.User
	err := row.Scan(&user.ID, &user.Email)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (p *PostgresRepository) Close() error {
	return p.db.Close()
}