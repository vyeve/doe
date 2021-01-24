package data

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

type sqlImpl struct {
	*sql.DB
}

type (
	SQLTx interface {
		Exec(query string, args ...interface{}) (sql.Result, error)
		ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
		Query(query string, args ...interface{}) (*sql.Rows, error)
		QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
		QueryRow(query string, args ...interface{}) *sql.Row
		QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
	}

	SQLDb interface {
		SQLTx
		Close() error
		DBConn() *sql.DB
	}
)

func NewSQLTx() (SQLDb, error) {
	user := os.Getenv("POSTGRES_USER")
	dbName := os.Getenv("POSTGRES_DB")
	port := os.Getenv("POSTGRES_PORT")
	host := os.Getenv("POSTGRES_HOST")
	connStr := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable",
		host, port, user, dbName)

	conn, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	return &sqlImpl{conn}, nil
}

func (s *sqlImpl) DBConn() *sql.DB {
	return s.DB
}

func ExecuteTx(ctx context.Context, db *sql.DB, fn func(tx SQLTx) error) (err error) {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		if err == nil {
			err = tx.Commit()
		} else {
			_ = tx.Rollback()
		}
	}()
	return fn(tx)
}
