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
	return 0, nil
}

// Return snippet
func (m *SnippetModel) Get(id int) (*Snippet, error) {
	return nil, nil
}

// Return 10 recent snippets
func (m *SnippetModel) Latest() ([]*Snippet, error) {
	return nil, nil
}
