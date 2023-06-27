package dictionary

import (
	"database/sql"
	"errors"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type Entry struct {
	Definition string
	Date       time.Time
}

func (e Entry) String() string {
	return e.Definition + " (added on " + e.Date.Format("2 Jan 2006 15:04") + ")"
}

type Dictionary struct {
	db *sql.DB
}

func New() *Dictionary {
	db, err := sql.Open("sqlite3", "./dictionary.db")
	if err != nil {
		panic(err)
	}

	statement, err := db.Prepare("CREATE TABLE IF NOT EXISTS dictionary (word TEXT, definition TEXT, date TIMESTAMP)")
	if err != nil {
		panic(err)
	}
	_, err = statement.Exec()
	if err != nil {
		panic(err)
	}

	return &Dictionary{
		db: db,
	}
}

func (d *Dictionary) Add(word string, definition string) error {
	tx, err := d.db.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare("INSERT INTO dictionary(word, definition, date) values(?,?,?)")
	if err != nil {
		return err
	}

	_, err = stmt.Exec(word, definition, time.Now())
	if err != nil {
		return err
	}

	tx.Commit()

	return nil
}

func (d *Dictionary) Get(word string) (Entry, error) {
	row := d.db.QueryRow("SELECT word, definition, date FROM dictionary WHERE word = ?", word)

	var e Entry
	if err := row.Scan(&word, &e.Definition, &e.Date); err != nil {
		if err == sql.ErrNoRows {
			return Entry{}, errors.New("word not found")
		}
		return Entry{}, err
	}
	return e, nil
}

func (d *Dictionary) Remove(word string) error {
	stmt, err := d.db.Prepare("DELETE FROM dictionary WHERE word=?")
	if err != nil {
		return err
	}

	_, err = stmt.Exec(word)
	if err != nil {
		return err
	}

	return nil
}

func (d *Dictionary) List() ([]string, map[string]Entry, error) {
	rows, err := d.db.Query("SELECT word, definition, date FROM dictionary")
	if err != nil {
		return nil, nil, err
	}

	words := make([]string, 0)
	entries := make(map[string]Entry)

	for rows.Next() {
		var word string
		var e Entry
		if err := rows.Scan(&word, &e.Definition, &e.Date); err != nil {
			return nil, nil, err
		}
		words = append(words, word)
		entries[word] = e
	}
	return words, entries, nil
}
