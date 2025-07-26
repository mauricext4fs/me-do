package repository

import (
	"database/sql"
	"errors"
	"log"
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

	CREATE TABLE IF NOT EXISTS task_positions (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		task_id INTEGER NOT NULL,
		position INTEGER DEFAULT 0,
		label TEXT DEFAULT 'TODO'
	);


	CREATE TABLE IF NOT EXISTS task_timers (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		task_id INTEGER NOT NULL,
		start_timestamp INTEGER DEFAULT 0,
		end_timestamp INTEGER DEFAULT 0
	);
	

INSERT INTO
tasks
	(id, position, title)
SELECT
	1, 1, "Sample task"
WHERE NOT EXISTS(
	SELECT 1 FROM tasks WHERE id = 1
);
 
INSERT INTO
task_positions
	(id, task_id, position, label)
SELECT
	1, 1, 1, "TODO"
WHERE NOT EXISTS(
	SELECT 1 FROM task_positions WHERE id = 1
);
	
`
	_, err := repo.Conn.Exec(query)
	return err
}

func (repo *SQLiteRepository) InsertTask(tasks Tasks) (*Tasks, error) {
	stmt := "INSERT INTO tasks (position, title, priority, created_at, updated_at) VALUES ((SELECT MAX(position) +1 FROM tasks), ?, ?, ?, ?)"

	res, err := repo.Conn.Exec(stmt, tasks.Title, tasks.Priority, time.Now().Unix(), time.Now().Unix())
	if err != nil {
		return nil, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}

	tasks.ID = id
	_ = repo.PushPosition()
	// Add to position table
	ps := &Positions{
		TaskID: id,
		Label:  "TODO",
	}
	repo.InsertPosition(*ps)

	return &tasks, nil
}

func (repo *SQLiteRepository) InsertPosition(position Positions) (*Positions, error) {
	query := `
	INSERT INTO 
		task_positions
		(task_id, position, label)
	VALUES 
		(?, 1, ?);
	`
	res, err := repo.Conn.Exec(query, position.TaskID, position.Label)
	if err != nil {
		return nil, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}

	position.ID = id

	return &position, nil
}

func (repo *SQLiteRepository) PushPosition() error {
	query := `
	UPDATE 
		task_positions
	SET position = position +1;
	`
	res, err := repo.Conn.Exec(query)

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

func (repo *SQLiteRepository) AllTODOTasks() ([]Tasks, error) {
	query := `
		SELECT
			t.id, t.title, IFNULL(tp.position, 9999999) AS pos, t.status, t.priority, t.created_at, t.updated_at
		FROM
			tasks t
		LEFT JOIN
			task_positions tp ON (t.id = tp.task_id AND label = 'TODO')
		WHERE
			t.status != 'Done'
		ORDER BY
			pos ASC, t.id DESC
	`
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

func (repo *SQLiteRepository) DownPosition(id int64, curPos int64, label string) error {
	if id == 0 {
		return errors.New("Invalid Updated ID")
	}

	// First we need to "downgrade" the task that is in the new position
	stmt := `
	UPDATE
		task_positions
	SET
		position = position-1
	WHERE 
		position = ?
		AND label = ?
	`
	log.Println("Setting existing task with new Position: ", curPos)
	res, err := repo.Conn.Exec(stmt, curPos, label)
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

	// Then we up the task position
	stmt = `
	UPDATE
		task_positions
	SET
		position = ?
	WHERE
		task_id = ?
		AND label = ?
	`

	log.Println("Setting task Position: ", curPos+1, " for task_id : ", id)
	res, err = repo.Conn.Exec(stmt, (curPos + 1), id, label)

	if err != nil {
		return err
	}

	rowAffected, err = res.RowsAffected()
	if err != nil {
		return err
	}

	if rowAffected == 0 {
		return errUpdateFailed
	}

	return nil

}

func (repo *SQLiteRepository) UpPosition(id int64, curPos int64, label string) error {
	if id == 0 {
		return errors.New("Invalid Updated ID")
	}

	// First we need to "Upgrade" the task that is in the new position
	stmt := `
	UPDATE
		task_positions
	SET
		position = position+1
	WHERE 
		position = ?
		AND label = ?
	`
	log.Println("Setting lower task Position to: ", curPos+1)
	res, err := repo.Conn.Exec(stmt, curPos, label)
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

	// Then we up the task position
	stmt = `
	UPDATE
		task_positions
	SET
		position = ?
	WHERE
		task_id = ?
		AND label = ?
	`
	log.Println("Setting task Position: ", curPos-1, " for task_id : ", id)
	res, err = repo.Conn.Exec(stmt, curPos-1, id, label)

	if err != nil {
		return err
	}

	rowAffected, err = res.RowsAffected()
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
		return errors.New("Invalid Task ID")
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

func (repo *SQLiteRepository) UpdatePriority(id int64, status string) error {
	if id == 0 {
		return errors.New("Invalid Task ID")
	}

	stmt := "UPDATE tasks SET priority = ?, updated_at = ? WHERE id = ?"
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
		return errors.New("Invalid Task ID")
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
