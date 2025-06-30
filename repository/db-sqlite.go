package repository

import (
	"database/sql"
	"errors"
	"time"
)

type SQLiteRepository struct {
	Conn *sql.DB
}

func NewSQLiteRepository(db *sql.DB) *SQLiteRepository {
	return &SQLiteRepository{
		Conn: db,
	}
}

func (repo *SQLiteRepository) Migrate() error {
	query := `
	
	CREATE TABLE IF NOT EXISTS tasks (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		position INTEGER DEFAULT 0,
		title TEXT NOT NULL,
		status TEXT DEFAULT 'Not started',
		priority TEXT DEFAULT '',
		created_at INTEGER DEFAULT 0,
		updated_at INTEGER DEFAULT 0
	);


	CREATE TABLE IF NOT EXISTS task_timers (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		task_id INTEGER NOT NULL,
		start_timestamp INTEGER DEFAULT 0,
		end_timestamp INTEGER DEFAULT 0
	);
	

	CREATE TABLE IF NOT EXISTS labels (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		position INTEGER DEFAULT 0,
		title TEXT NOT NULL,
		color TEXT NOT NULL
	);
	
	CREATE TABLE IF NOT EXISTS task_labels (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		position INTEGER DEFAULT 0,
		task_id INTEGER NOT NULL
	);
	
`
	_, err := repo.Conn.Exec(query)
	return err
}

func (repo *SQLiteRepository) InsertTask(tasks Tasks) (*Tasks, error) {
	stmt := "INSERT INTO Tasks (position, title, priority, created_at, updated_at) values (MAX(position), ?, ?, ?, ?)"

	res, err := repo.Conn.Exec(stmt, tasks.Title, tasks.Priority, time.Now().Unix(), time.Now().Unix())
	if err != nil {
		return nil, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}

	tasks.ID = id

	return &tasks, nil
}

func (repo *SQLiteRepository) AllTasks() ([]Tasks, error) {
	query := "SELECT id, title, position, status, priority, created_at, updated_at FROM tasks ORDER BY position DESC, id DESC"
	rows, err := repo.Conn.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var all []Tasks
	for rows.Next() {
		var a Tasks
		var cA int64
		var uA int64
		err := rows.Scan(
			&a.ID,
			&a.Title,
			&a.Position,
			&a.Status,
			&a.Priority,
			&cA,
			&uA,
		)
		if err != nil {
			return nil, err
		}
		a.CreatedAt = time.Unix(cA, 0)
		a.UpdatedAt = time.Unix(uA, 0)
		all = append(all, a)
	}

	return all, nil
}

func (repo *SQLiteRepository) GetTaskByID(id int) (*Tasks, error) {
	row := repo.Conn.QueryRow("SELECT id, title, created_at, updated_at FROM tasks WHERE id = ?", id)

	var a Tasks
	var startUnixTime int64
	var endUnixTime int64
	err := row.Scan(
		&a.ID,
		&a.Title,
		&startUnixTime,
		&endUnixTime,
	)

	if err != nil {
		return nil, err
	}

	a.CreatedAt = time.Unix(startUnixTime, 0)
	a.UpdatedAt = time.Unix(endUnixTime, 0)

	return &a, nil
}

func (repo *SQLiteRepository) UpdatePosition(id int64, newPos int64) error {
	if id == 0 {
		return errors.New("Invalid Updated ID")
	}

	stmt := "UPDATE tasks SET position = ?, updated_at = ? WHERE id = ?"
	res, err := repo.Conn.Exec(stmt, newPos, time.Now().Unix(), id)

	if err != nil {
		return err
	}

	rowAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowAffected == 0 {
		return errUpdateFailed
	}

	return nil

}

func (repo *SQLiteRepository) UpdateStatus(id int64, status string) error {
	if id == 0 {
		return errors.New("Invalid Updated ID")
	}

	stmt := "UPDATE tasks SET status = ?, updated_at = ? WHERE id = ?"
	res, err := repo.Conn.Exec(stmt, status, time.Now().Unix(), id)

	if err != nil {
		return err
	}

	rowAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowAffected == 0 {
		return errUpdateFailed
	}

	return nil

}

func (repo *SQLiteRepository) UpdateTask(id int64, updated Tasks) error {
	if id == 0 {
		return errors.New("Invalid Updated ID")
	}

	stmt := "UPDATE tasks SET updated_at = ? WHERE id = ?"
	res, err := repo.Conn.Exec(stmt, updated.UpdatedAt.Unix(), id)

	if err != nil {
		return err
	}

	rowAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowAffected == 0 {
		return errUpdateFailed
	}

	return nil

}

func (repo *SQLiteRepository) DeleteTask(id int64) error {
	res, err := repo.Conn.Exec("DELETE FROM tasks WHERE id = ?", id)
	if err != nil {
		return err

	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errDeleteFailed
	}

	return nil

}
