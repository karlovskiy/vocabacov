package database

import (
	"database/sql"
	"embed"
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"io/fs"
	"log/slog"
	"os"

	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	_ "github.com/mattn/go-sqlite3"

	"github.com/karlovskiy/vocabacov/internal/translate"
)

const (
	EnvDatabase = "VOCABACOV_DB_PATH"
)

//go:embed migrations/*.sql
var migrationsFs embed.FS

func OpenDb(withMigrate bool) (*sql.DB, error) {
	dbPath := os.Getenv(EnvDatabase)
	if dbPath == "" {
		return nil, fmt.Errorf("database connection string not found in environment variable %s", EnvDatabase)
	}
	slog.Info("open db", "path", dbPath)
	_, err := os.Stat(dbPath)
	if errors.Is(err, fs.ErrNotExist) {
		return nil, fmt.Errorf("db file not exist: %w", err)
	}
	if err != nil {
		return nil, fmt.Errorf("db file error: %w", err)
	}
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("db open error: %w", err)
	}
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("db ping error: %w", err)
	}
	if withMigrate {
		driver, err := sqlite3.WithInstance(db, &sqlite3.Config{})
		if err != nil {
			return nil, fmt.Errorf("migration db driver error: %w", err)
		}
		d, err := iofs.New(migrationsFs, "migrations")
		if err != nil {
			return nil, fmt.Errorf("migration source driver error: %w", err)
		}
		m, err := migrate.NewWithInstance("iofs", d, "vocabacov", driver)
		if err != nil {
			return nil, fmt.Errorf("migration creation error: %w", err)
		}
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			return nil, fmt.Errorf("migration error: %w", err)
		}
	}
	return db, nil
}

func SavePhrase(db *sql.DB, phrase *translate.Phrase) error {
	res, err := db.Exec("INSERT INTO PHRASES(LANG, PHRASE, TRANSLATION) VALUES(?, ?, ?)",
		phrase.Lang, phrase.Phrase, phrase.Translation)
	if err != nil {
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	slog.Info("inserted phrase", "rows", rows, "phrase", phrase)
	return nil
}

func SetPhrasesStatus(db *sql.DB, lang, status string) error {
	res, err := db.Exec("UPDATE PHRASES SET STATUS=? WHERE LANG=?", status, lang)
	if err != nil {
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	slog.Info("set phrases status", "lang", lang, "status", status, "rows", rows)
	return nil
}

func LoadActivePhrases(db *sql.DB, lang string) ([]translate.Phrase, error) {
	rows, err := db.Query("SELECT LANG, PHRASE, TRANSLATION FROM PHRASES WHERE LANG=? AND STATUS='ACTIVE' ORDER BY ID", lang)
	if err != nil {
		return nil, fmt.Errorf("error loading %s phrases: %v", lang, err)
	}
	defer rows.Close()
	phrases := make([]translate.Phrase, 0, 100)
	for rows.Next() {
		var p translate.Phrase
		if err := rows.Scan(&p.Lang, &p.Phrase, &p.Translation); err != nil {
			return nil, fmt.Errorf("error scan phrases")
		}
		phrases = append(phrases, p)
	}
	return phrases, nil
}

func LoadAllPhrases(db *sql.DB) ([]translate.Phrase, error) {
	rows, err := db.Query("SELECT LANG, PHRASE, TRANSLATION FROM PHRASES ORDER BY ID")
	if err != nil {
		return nil, fmt.Errorf("error loading all phrases: %v", err)
	}
	defer rows.Close()
	phrases := make([]translate.Phrase, 0, 100)
	for rows.Next() {
		var p translate.Phrase
		if err := rows.Scan(&p.Lang, &p.Phrase, &p.Translation); err != nil {
			return nil, fmt.Errorf("error scan phrases")
		}
		phrases = append(phrases, p)
	}
	return phrases, nil
}
