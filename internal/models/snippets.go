package models

import (
	"database/sql"
	"time"
)

type Snippet struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}

// Wrapper for sql connection pool
type SnippetModel struct {
	DB *sql.DB
}

// Insert new snippet
func (m *SnippetModel) Insert(title string, content string, expires int) (int, error) {
	stmt := `INSERT INTO snippets (title, content, created, expires)
			VALUES(?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))`
	result, err := m.DB.Exec(stmt, title, content, expires)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil

}

// Return snippet
func (m *SnippetModel) Get(id int) (*Snippet, error) {
	return nil, nil
}

// Return 10 recent snippets
func (m *SnippetModel) Latest() ([]*Snippet, error) {
	return nil, nil
}
