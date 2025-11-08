package store

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

// Store manages the database connection and queries
type Store struct {
	DB *sql.DB
}

// NewStore creates a new Store
func NewStore(dbName string) (*Store, error) {
	db, err := sql.Open("sqlite3", dbName)
	if err != nil {
		return nil, err
	}
	createTable := `CREATE TABLE IF NOT EXISTS flashcards (
			  id INTEGER PRIMARY KEY AUTOINCREMENT,
			  file TEXT,
			  question TEXT,
			  answer TEXT,
			  revisitin INTEGER DEFAULT 0
		  );`
	_, err = db.Exec(createTable)
	if err != nil {
		return nil, err
	}
	return &Store{DB: db}, nil
}

// GetFlashcardsForReview returns all flashcards that are due for review
func (s *Store) GetFlashcardsForReview() ([]Flashcard, error) {
	rows, err := s.DB.Query("SELECT id, question, answer, revisitin FROM flashcards WHERE revisitin <= 0 ORDER BY id ASC")
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = rows.Close()
	}()

	var flashcards []Flashcard
	for rows.Next() {
		var fc Flashcard
		err := rows.Scan(&fc.ID, &fc.Question, &fc.Answer, &fc.RevisitIn)
		if err != nil {
			return nil, err
		}
		flashcards = append(flashcards, fc)
	}
	return flashcards, nil
}

// GetAllFlashcards returns all flashcards ordered by revisitin ascending
func (s *Store) GetAllFlashcards() ([]Flashcard, error) {
	rows, err := s.DB.Query("SELECT id, file, question, answer, revisitin FROM flashcards ORDER BY revisitin ASC, id ASC")
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = rows.Close()
	}()
	var flashcards []Flashcard
	for rows.Next() {
		var fc Flashcard
		if err := rows.Scan(&fc.ID, &fc.File, &fc.Question, &fc.Answer, &fc.RevisitIn); err != nil {
			return nil, err
		}
		flashcards = append(flashcards, fc)
	}
	return flashcards, nil
}

// DeleteFlashcard deletes a flashcard by id
func (s *Store) DeleteFlashcard(id int) error {
	_, err := s.DB.Exec("DELETE FROM flashcards WHERE id=?", id)
	return err
}

// UpdateFlashcardFull updates all editable fields of a flashcard
func (s *Store) UpdateFlashcardFull(fc Flashcard) error {
	_, err := s.DB.Exec("UPDATE flashcards SET file=?, question=?, answer=?, revisitin=? WHERE id=?", fc.File, fc.Question, fc.Answer, fc.RevisitIn, fc.ID)
	return err
}

// UpdateFlashcard updates a flashcard's revisitin date
func (s *Store) UpdateFlashcard(fc Flashcard) error {
	_, err := s.DB.Exec("UPDATE flashcards SET revisitin=? WHERE id=?", fc.RevisitIn, fc.ID)
	return err
}

// IsFileProcessed checks if a file has already been processed
func (s *Store) IsFileProcessed(filePath string) (bool, error) {
	var count int
	row := s.DB.QueryRow("SELECT COUNT(*) FROM flashcards WHERE file = ?", filePath)
	err := row.Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// InsertFlashcard inserts a new flashcard into the database
func (s *Store) InsertFlashcard(fc Flashcard) error {
	_, err := s.DB.Exec("INSERT INTO flashcards (file, question, answer, revisitin) VALUES (?, ?, ?, ?)", fc.File, fc.Question, fc.Answer, fc.RevisitIn)
	return err
}

// Close closes the database connection
func (s *Store) Close() {
	_ = s.DB.Close()
}
