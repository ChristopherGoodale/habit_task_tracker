package task

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

var ErrNotFound = errors.New("task not found")

// Store is a thin, explicit wrapper around database/sql — every query is
// spelled out here rather than generated, so the SQL that runs is always the
// SQL you can see.
type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

func (s *Store) List(ctx context.Context) ([]Task, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT id, title, description, done, created_at, updated_at
		FROM tasks
		ORDER BY created_at DESC`)
	if err != nil {
		return nil, fmt.Errorf("list tasks: %w", err)
	}
	defer rows.Close()

	tasks := []Task{}
	for rows.Next() {
		var t Task
		if err := rows.Scan(&t.ID, &t.Title, &t.Description, &t.Done, &t.CreatedAt, &t.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan task: %w", err)
		}
		tasks = append(tasks, t)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("list tasks: %w", err)
	}

	return tasks, nil
}

func (s *Store) Create(ctx context.Context, title, description string) (Task, error) {
	var t Task
	err := s.db.QueryRowContext(ctx, `
		INSERT INTO tasks (title, description)
		VALUES ($1, $2)
		RETURNING id, title, description, done, created_at, updated_at`,
		title, description,
	).Scan(&t.ID, &t.Title, &t.Description, &t.Done, &t.CreatedAt, &t.UpdatedAt)
	if err != nil {
		return Task{}, fmt.Errorf("create task: %w", err)
	}
	return t, nil
}

func (s *Store) Update(ctx context.Context, id int64, title, description string, done bool) (Task, error) {
	var t Task
	err := s.db.QueryRowContext(ctx, `
		UPDATE tasks
		SET title = $2, description = $3, done = $4, updated_at = now()
		WHERE id = $1
		RETURNING id, title, description, done, created_at, updated_at`,
		id, title, description, done,
	).Scan(&t.ID, &t.Title, &t.Description, &t.Done, &t.CreatedAt, &t.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return Task{}, ErrNotFound
	}
	if err != nil {
		return Task{}, fmt.Errorf("update task: %w", err)
	}
	return t, nil
}

func (s *Store) Delete(ctx context.Context, id int64) error {
	result, err := s.db.ExecContext(ctx, `DELETE FROM tasks WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("delete task: %w", err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("delete task: %w", err)
	}
	if rows == 0 {
		return ErrNotFound
	}
	return nil
}
