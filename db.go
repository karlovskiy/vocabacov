package vocabacov

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"io/fs"
	"log"
	"os"
)

func NewDb() (*sql.DB, error) {
	dbPath := os.Getenv(EnvDatabase)
	if dbPath == "" {
		return nil, fmt.Errorf("database connection string not found in environment variable %s", EnvDatabase)
	}
	log.Printf("db path: %q\n", dbPath)
	_, err := os.Stat(dbPath)
	if errors.Is(err, fs.ErrNotExist) {
		return nil, fmt.Errorf("db file not exist: %w", err)
	}
	if err != nil {
		return nil, fmt.Errorf("db file error: %w", err)
	}
	db, err := sql.Open("sqlite3", dbPath)
	if err := db.Ping(); err != nil {
		return nil, err
	}
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS PHRASES(ID INTEGER PRIMARY KEY, LANG TEXT NOT NULL, PHRASE TEXT NOT NULL)")
	if err != nil {
		return db, err
	}
	return db, err
}

func savePhrase(db *sql.DB, lang, phrase string) error {
	res, err := db.Exec("INSERT INTO PHRASES(LANG, PHRASE) VALUES(?, ?)", lang, phrase)
	if err != nil {
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	log.Printf("inserted phrase %q, lang: %s, rows: %d", phrase, lang, rows)
	return nil
}
