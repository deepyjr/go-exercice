package dictionary

import (
	"errors"
	"time"
)

type Entry struct {
	Definition string
	Date       time.Time
}

func (e Entry) String() string {
	return e.Definition + " (added on " + e.Date.String() + ")"
}

type Dictionary struct {
	entries map[string]Entry
}

func New() *Dictionary {
	return &Dictionary{
		entries: make(map[string]Entry),
	}
}

func (d *Dictionary) Add(word string, definition string) {
	d.entries[word] = Entry{
		Definition: definition,
		Date:       time.Now(),
	}
}

func (d *Dictionary) Get(word string) (Entry, error) {
	entry, exists := d.entries[word]
	if !exists {
		return Entry{}, errors.New("word not found")
	}
	return entry, nil
}

func (d *Dictionary) Remove(word string) {
	delete(d.entries, word)
}

func (d *Dictionary) List() ([]string, map[string]Entry) {
	words := make([]string, 0, len(d.entries))
	for word := range d.entries {
		words = append(words, word)
	}
	return words, d.entries
}
