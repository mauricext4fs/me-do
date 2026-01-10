package repository

func (repo *SQLiteRepository) GetDBSchema() string {
	schema := `
	
	CREATE TABLE IF NOT EXISTS tasks (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		position INTEGER DEFAULT 0,
		title TEXT NOT NULL,
		status TEXT DEFAULT 'Not started',
		priority TEXT DEFAULT 'Low',
		created_at INTEGER DEFAULT 0,
		created_by INTEGER DEFAULT 1,
		updated_at INTEGER DEFAULT 0,
		updated_by INTEGER DEFAULT 1
	);

	CREATE TABLE IF NOT EXISTS task_positions (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		task_id INTEGER NOT NULL,
		position INTEGER NOT NULL,
		label TEXT DEFAULT 'TODO'
	);

	CREATE TABLE IF NOT EXISTS task_notes (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		task_id INTEGER NOT NULL,
		note TEXT NOT NULL,
		created_at INTEGER DEFAULT 0,
		created_by INTEGER DEFAULT 1,
		updated_at INTEGER DEFAULT 0,
		updated_by INTEGER DEFAULT 1
	);

	CREATE TABLE IF NOT EXISTS files (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		filetype TEXT NOT NULL,
		deleted INTEGER DEFAULT 0,
		created_at INTEGER DEFAULT 0,
		created_by INTEGER DEFAULT 1
	);

	CREATE TABLE IF NOT EXISTS task_note_files (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		note_id INTEGER NOT NULL,
		file_id INTEGER NOT NULL
	);

	CREATE TABLE IF NOT EXISTS task_timers (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		task_id INTEGER NOT NULL,
		start_timestamp INTEGER DEFAULT 0,
		end_timestamp INTEGER DEFAULT 0,
		created_by INTEGER DEFAULT 1
	);

`
	return schema

}

func (repo *SQLiteRepository) GetDefaultData() string {

	query := `

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

	return query
}
