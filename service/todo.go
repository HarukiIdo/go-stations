package service

import (
	"context"
	"database/sql"
	"time"

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
	ins, err := s.db.PrepareContext(ctx, insert)
	if err != nil {
		return nil, err
	}

	// ステートメントを実行するが、行は返さない
	res, err := ins.ExecContext(ctx, subject, description)
	if err != nil {
		return nil, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}

	conf, err := s.db.PrepareContext(ctx, confirm)
	if err != nil {
		return nil, err
	}

	todo := &model.TODO{
		ID: id,
	}

	// ステートメントを実行し、行を返す
	err = conf.QueryRowContext(ctx, id).Scan(&todo.Subject, &todo.Description, &todo.CreatedAt, &todo.UpdatedAt)
	if err != nil {
		return nil, err
	}
	//log.Println(todo)

	return todo, err
}

// ReadTODO reads TODOs on DB.
func (s *TODOService) ReadTODO(ctx context.Context, prevID, size int64) ([]*model.TODO, error) {
	const (
		read       = `SELECT id, subject, description, created_at, updated_at FROM todos ORDER BY id DESC LIMIT ?`
		readWithID = `SELECT id, subject, description, created_at, updated_at FROM todos WHERE id < ? ORDER BY id DESC LIMIT ?`
	)

	readStm, err := s.db.PrepareContext(ctx, read)
	if err != nil {
		return nil, err
	}

	readwithIDStm, err := s.db.PrepareContext(ctx, readWithID)
	if err != nil {
		return nil, err
	}

	var rows *sql.Rows

	// PrevIDを持っているかによりクエリを変更
	if prevID == 0 {
		rows, err = readStm.QueryContext(ctx, size)
	} else {
		rows, err = readwithIDStm.QueryContext(ctx, prevID, size)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	todos := []*model.TODO{}
	for rows.Next() {
		todo := &model.TODO{}
		if err := rows.Scan(&todo.ID, &todo.Subject, &todo.Description, &todo.CreatedAt, &todo.UpdatedAt); err != nil {
			return nil, err
		}
		todos = append(todos, todo)
	}

	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return todos, err
}

// UpdateTODO updates the TODO on DB.
func (s *TODOService) UpdateTODO(ctx context.Context, id int64, subject, description string) (*model.TODO, error) {
	const (
		update  = `UPDATE todos SET subject = ?, description = ? WHERE id = ?`
		confirm = `SELECT subject, description, created_at, updated_at FROM todos WHERE id = ?`
	)

	updateStmt, err := s.db.PrepareContext(ctx, update)
	if err != nil {
		return nil, err
	}

	res, err := updateStmt.ExecContext(ctx, subject, description, id)
	if err != nil {
		return nil, err
	}

	num, err := res.RowsAffected()
	if num == 0 {
		errNotFound := &model.ErrNotFound{
			When: time.Now(),
			What: "todos not found",
		}
		return nil, errNotFound
	}

	conf, err := s.db.PrepareContext(ctx, confirm)
	if err != nil {
		return nil, err
	}

	todo := &model.TODO{
		ID:          id,
		Subject:     subject,
		Description: description,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	err = conf.QueryRowContext(ctx, id).Scan(&todo.Subject, &todo.Description, &todo.CreatedAt, &todo.UpdatedAt, &todo.ID)

	return todo, nil
}

// DeleteTODO deletes TODOs on DB by ids.
func (s *TODOService) DeleteTODO(ctx context.Context, ids []int64) error {
	const deleteFmt = `DELETE FROM todos WHERE id IN (?%s)`

	return nil
}
