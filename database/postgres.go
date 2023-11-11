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
	_, err := p.db.ExecContext(ctx, "INSERT INTO users (id, email, password) VALUES ($1, $2, $3)",user.ID, user.Email, user.Password)
	return err
}


func (p *PostgresRepository) InsertPost(ctx context.Context, post *models.Post) error {
	_, err := p.db.ExecContext(ctx, "INSERT INTO posts (id, post_content, user_id) VALUES ($1, $2, $3)",
	post.ID, post.PostContent, post.UserId)
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

func (p *PostgresRepository) GetPostById(ctx context.Context, id string) (*models.Post, error) {
	row := p.db.QueryRowContext(ctx, "SELECT id, post_content, user_id, created_at FROM posts WHERE id = $1", id)
	var post models.Post
	err := row.Scan(&post.ID, &post.PostContent, &post.UserId, &post.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &post, nil
}

func (p *PostgresRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	row := p.db.QueryRowContext(ctx, "SELECT id, email, password FROM users WHERE email = $1", email)
	var user models.User
	err := row.Scan(&user.ID, &user.Email, &user.Password)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (p *PostgresRepository) UpdatePost(ctx context.Context, post *models.Post) error {
	_, err := p.db.ExecContext(ctx, "UPDATE posts SET post_content = $1 WHERE id = $2 and user_id = $3", post.PostContent, post.ID, post.UserId)
	return err
}

func (p *PostgresRepository) DeletePost(ctx context.Context, id string, userid int64) error {
	_, err := p.db.ExecContext(ctx, "DELETE FROM posts WHERE id = $1 and user_id = $2", id, userid)
	return err
}



func (p *PostgresRepository) Close() error {
	return p.db.Close()
}