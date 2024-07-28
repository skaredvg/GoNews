package postgres

import (
	"context"
	"log"
	"skillfactory/GoNews/pkg/storage"

	"github.com/jackc/pgx/v4/pgxpool"
)

type Storage struct {
	db *pgxpool.Pool
}

func New(connstr string) (*Storage, error) {
	db, err := pgxpool.Connect(context.Background(), connstr)
	if err != nil {
		return nil, err
	}
	return &Storage{
		db: db,
	}, nil
}

func (s *Storage) Posts() ([]storage.Post, error) {
	ps := make([]storage.Post, 0)
	sql := `SELECT p.id, p.title, p.content, p.author_id, a.name, p.created_at, p.published_at FROM posts.public.posts p 
			JOIN posts.public.authors a on a.id = p.author_id`
	rows, err := s.db.Query(context.Background(), sql)
	if err != nil {
		return ps, err
	}
	defer rows.Close()

	for rows.Next() {
		p := storage.Post{}
		if err := rows.Scan(&p.ID, &p.Title, &p.Content, &p.AuthorID, &p.AuthorName, &p.CreatedAt, &p.PublishedAt); err != nil {
			return ps, err
		}
		ps = append(ps, p)
	}
	log.Println(ps)
	return ps, nil
}

func (s *Storage) AddPost(p storage.Post) error {
	sql := `INSERT INTO posts.public.posts (title, content, author_id, created_at, published_at)
			VALUES($1, $2, $3, $4, $5) RETURNING id`
	row := s.db.QueryRow(context.Background(), sql, p.Title, p.Content, p.AuthorID, p.CreatedAt, p.PublishedAt)

	return row.Scan(&p.ID)
}

func (s *Storage) UpdatePost(p storage.Post) error {
	sql := `UPDATE posts.public.posts
			SET title = $1,
				content = $2,
				author_id = $3,
				created_at = $4,
				published_at = $5
			WHERE id = $6`
	_, err := s.db.Exec(context.Background(), sql, p.Title, p.Content, p.AuthorID, p.CreatedAt, p.PublishedAt, p.ID)
	return err

}

func (s *Storage) DeletePost(p storage.Post) error {
	sql := `DELETE FROM posts.public.posts p
			WHERE p.id = $1`
	_, err := s.db.Exec(context.Background(), sql, p.ID)

	return err
}
