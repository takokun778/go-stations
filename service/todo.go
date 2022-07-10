package service

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/TechBowl-japan/go-stations/model"
)

// A TODOService implements CRUD of TODO entities.
type TODOService struct {
	db *sql.DB
}

// NewTODOService returns new TODOService.
func NewTODOService(db *sql.DB) *TODOService {
	return &TODOService{
		db: db,
	}
}

// CreateTODO creates a TODO on DB.
func (s *TODOService) CreateTODO(ctx context.Context, subject, description string) (*model.TODO, error) {
	const (
		insert  = `INSERT INTO todos(subject, description) VALUES(?, ?)`
		confirm = `SELECT subject, description, created_at, updated_at FROM todos WHERE id = ?`
	)

	res, err := s.db.ExecContext(ctx, insert, subject, description)

	if err != nil {
		return nil, err
	}

	id, err := res.LastInsertId()

	if err != nil {
		return nil, err
	}

	row := s.db.QueryRowContext(ctx, confirm, id)

	var todo model.TODO

	if err := row.Scan(&todo.Subject, &todo.Description, &todo.CreatedAt, &todo.UpdatedAt); err != nil {
		return nil, err
	}

	todo.ID = id

	return &todo, nil
}

// ReadTODO reads TODOs on DB.
func (s *TODOService) ReadTODO(ctx context.Context, prevID, size int) ([]*model.TODO, error) {
	const (
		read       = `SELECT id, subject, description, created_at, updated_at FROM todos ORDER BY id DESC LIMIT ?`
		readWithID = `SELECT id, subject, description, created_at, updated_at FROM todos WHERE id < ? ORDER BY id DESC LIMIT ?`
	)

	// この実装はどうなのだろうか...
	if size == 0 {
		size = 100
	}

	if prevID == 0 {
		rows, err := s.db.QueryContext(ctx, read, size)

		if err != nil {
			return nil, err
		}

		todos := make([]*model.TODO, 0, size)

		for rows.Next() {
			var todo model.TODO
			err := rows.Scan(&todo.ID, &todo.Subject, &todo.Description, &todo.CreatedAt, &todo.UpdatedAt)
			if err != nil {
				return nil, err
			}
			todos = append(todos, &todo)
		}

		return todos, nil
	}

	rows, err := s.db.QueryContext(ctx, readWithID, prevID, size)

	if err != nil {
		return nil, err
	}

	todos := make([]*model.TODO, 0, size)

	for rows.Next() {
		var todo model.TODO
		err := rows.Scan(&todo.ID, &todo.Subject, &todo.Description, &todo.CreatedAt, &todo.UpdatedAt)
		if err != nil {
			return nil, err
		}
		todos = append(todos, &todo)
	}

	return todos, nil
}

// UpdateTODO updates the TODO on DB.
func (s *TODOService) UpdateTODO(ctx context.Context, id int64, subject, description string) (*model.TODO, error) {
	const (
		update  = `UPDATE todos SET subject = ?, description = ? WHERE id = ?`
		confirm = `SELECT subject, description, created_at, updated_at FROM todos WHERE id = ?`
	)

	res, err := s.db.ExecContext(ctx, update, subject, description, id)

	if err != nil {
		return nil, err
	}

	count, err := res.RowsAffected()

	if err != nil {
		return nil, err
	}

	if count == 0 {
		return nil, model.ErrNotFound{}
	}

	row := s.db.QueryRowContext(ctx, confirm, id)

	var todo model.TODO

	if err := row.Scan(&todo.Subject, &todo.Description, &todo.CreatedAt, &todo.UpdatedAt); err != nil {
		return nil, err
	}

	todo.ID = id

	return &todo, nil
}

// DeleteTODO deletes TODOs on DB by ids.
func (s *TODOService) DeleteTODO(ctx context.Context, ids []int64) error {
	const deleteFmt = `DELETE FROM todos WHERE id IN (?%s)`

	q := make([]string, 0, len(ids))

	for range ids {
		q = append(q, "?")
	}

	sql := fmt.Sprintf("DELETE FROM todos WHERE id IN (%s)", strings.Join(q, ", "))

	if len(ids) == 0 {
		return nil
	}

	args := make([]interface{}, 0, len(ids))

	for _, id := range ids {
		args = append(args, id)
	}

	res, err := s.db.ExecContext(ctx, sql, args...)

	if err != nil {
		return err
	}

	count, err := res.RowsAffected()

	if err != nil {
		return err
	}

	if count == 0 {
		return model.ErrNotFound{}
	}

	return nil
}
