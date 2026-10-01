package history

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"errors"
	"modernc.org/sqlite"
	"os"
)

// backupForMigration creates a unique, validated online snapshot beside the
// current database. No schema migration currently calls this boundary.
func (s *Store) backupForMigration(ctx context.Context) (string, error) {
	if s == nil || s.db == nil || s.readOnly {
		return "", storeError(ErrMigrationFailed, errors.New("migration backup requires an open writer"))
	}
	if err := s.acquireWriter(ctx, ErrMigrationFailed); err != nil {
		return "", err
	}
	defer func() { <-s.writers }()
	var tokenBytes [16]byte
	if _, err := rand.Read(tokenBytes[:]); err != nil {
		return "", storeError(ErrMigrationFailed, err)
	}
	token := hex.EncodeToString(tokenBytes[:])
	temporary := s.path + ".backup-" + token + ".tmp"
	final := s.path + ".backup-" + token + ".sqlite3"
	file, err := os.OpenFile(temporary, os.O_CREATE|os.O_EXCL|os.O_RDWR, 0o600)
	if err != nil {
		return "", storeError(ErrMigrationFailed, err)
	}
	if err := file.Close(); err != nil {
		os.Remove(temporary)
		return "", storeError(ErrMigrationFailed, err)
	}
	committed := false
	defer func() {
		if !committed {
			os.Remove(temporary)
		}
	}()
	if err := secureFile(temporary, true, true); err != nil {
		return "", storeError(ErrMigrationFailed, err)
	}
	if err := secureSQLiteSidecars(temporary, false); err != nil {
		return "", storeError(ErrMigrationFailed, err)
	}
	source, err := s.db.Conn(ctx)
	if err != nil {
		return "", storeError(ErrMigrationFailed, err)
	}
	defer source.Close()
	err = source.Raw(func(driverConn any) error {
		backuper, ok := driverConn.(interface {
			NewBackup(string) (*sqlite.Backup, error)
		})
		if !ok {
			return errors.New("SQLite connection does not support online backup")
		}
		backup, err := backuper.NewBackup(dataSourceName(temporary, true))
		if err != nil {
			return err
		}
		finished := false
		defer func() {
			if !finished {
				backup.Finish()
			}
		}()
		for {
			if err := ctx.Err(); err != nil {
				return err
			}
			done, err := backup.Step(128)
			if err != nil {
				return err
			}
			if !done {
				break
			}
		}
		if err := backup.Finish(); err != nil {
			finished = true
			return err
		}
		finished = true
		return nil
	})
	if err != nil {
		return "", storeError(ErrMigrationFailed, err)
	}
	if err := secureFile(temporary, true, false); err != nil {
		return "", storeError(ErrMigrationFailed, err)
	}
	if err := secureSQLiteSidecars(temporary, true); err != nil {
		return "", storeError(ErrMigrationFailed, err)
	}
	if err := validateBackup(ctx, temporary); err != nil {
		return "", storeError(ErrMigrationFailed, err)
	}
	if err := os.Link(temporary, final); err != nil {
		return "", storeError(ErrMigrationFailed, err)
	}
	if err := secureFile(final, false, false); err != nil {
		os.Remove(final)
		return "", storeError(ErrMigrationFailed, err)
	}
	if err := os.Remove(temporary); err != nil {
		os.Remove(final)
		return "", storeError(ErrMigrationFailed, err)
	}
	committed = true
	return final, nil
}

func validateBackup(ctx context.Context, path string) error {
	db, err := sql.Open("sqlite", dataSourceName(path, false))
	if err != nil {
		return err
	}
	defer db.Close()
	db.SetMaxOpenConns(1)
	if err := db.PingContext(ctx); err != nil {
		return err
	}
	var version int
	if err := db.QueryRowContext(ctx, "PRAGMA user_version").Scan(&version); err != nil {
		return err
	}
	if version != 1 {
		return errors.New("backup schema version is not supported")
	}
	if err := verifySchema(ctx, db); err != nil {
		return err
	}
	var integrity string
	if err := db.QueryRowContext(ctx, "PRAGMA integrity_check").Scan(&integrity); err != nil || integrity != "ok" {
		return errors.New("online backup failed SQLite integrity validation")
	}
	tx, err := db.BeginTx(ctx, &sql.TxOptions{ReadOnly: true})
	if err != nil {
		return err
	}
	defer tx.Rollback()
	rows, err := tx.QueryContext(ctx, `SELECT history_entry_id FROM history_entries ORDER BY history_entry_id`)
	if err != nil {
		return err
	}
	ids := make([]string, 0, maxHistoryEntries)
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			rows.Close()
			return err
		}
		ids = append(ids, id)
	}
	if err := rows.Err(); err != nil {
		rows.Close()
		return err
	}
	if err := rows.Close(); err != nil {
		return err
	}
	for _, id := range ids {
		if _, _, err := readEntryByID(ctx, tx, id); err != nil {
			return err
		}
	}
	return tx.Commit()
}
