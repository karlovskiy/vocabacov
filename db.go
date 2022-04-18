package vocabacov

import (
	"database/sql"
	"embed"
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	_ "github.com/mattn/go-sqlite3"
	"io/fs"
	"log"
	"os"
)

//go:embed migrations/*.sql
var migrationsFs embed.FS

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
	driver, err := sqlite3.WithInstance(db, &sqlite3.Config{})
	if err != nil {
		return nil, err
	}
	d, err := iofs.New(migrationsFs, "migrations")
	if err != nil {
		log.Fatal(err)
	}
	m, err := migrate.NewWithInstance("iofs", d, "vocabacov", driver)
	if err != nil {
		return nil, err
	}
	if err := m.Up(); err != nil {
		return nil, err
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
