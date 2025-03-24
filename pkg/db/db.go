package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	migrations "github.com/fr-str/bingo/db"
	"github.com/fr-str/bingo/pkg/store"
	"github.com/fr-str/log"
	"github.com/fr-str/log/level"
	"github.com/pressly/goose/v3"
	"modernc.org/sqlite"
)

type db struct {
	w *sql.DB
}

func ConnectStore(ctx context.Context, filename string) (*store.Queries, error) {
	err := os.MkdirAll(filepath.Dir(filename), 0o755)
	if err != nil {
		return nil, err
	}
	w, err := sql.Open("sqlite", filename)
	if err != nil {
		return nil, err
	}

	go func() {
		<-ctx.Done()
		w.Close()
	}()

	// create tables
	goose.SetBaseFS(migrations.DBmigrations)
	goose.SetDialect("sqlite3")
	err = goose.UpContext(ctx, w, "migrations")
	if err != nil {
		panic(err)
	}

	_, err = w.ExecContext(ctx, "PRAGMA journal_mode=WAL")
	if err != nil {
		panic(err)
	}

	_, err = w.ExecContext(ctx, "PRAGMA synchronus=NORMAL")
	if err != nil {
		panic(err)
	}

	w.SetMaxOpenConns(1)
	w.SetConnMaxLifetime(0)
	w.SetConnMaxIdleTime(0)
	d := db{
		w: w,
	}
	return store.New(d), nil
}

func (s db) ExecContext(ctx context.Context, sql string, args ...any) (sql.Result, error) {
	ts := time.Now()
	res, err := s.w.ExecContext(ctx, sql, args...)
	logger(ctx, "ExecContext", sql, ts, args, res, err)
	return res, err
}

func (s db) PrepareContext(ctx context.Context, sql string) (*sql.Stmt, error) {
	ts := time.Now()
	stmt, err := s.w.PrepareContext(ctx, sql)
	logger(ctx, "PrepareContext", sql, ts, sql, nil, err)
	return stmt, err
}

func (s db) QueryContext(ctx context.Context, sql string, args ...any) (*sql.Rows, error) {
	ts := time.Now()
	rows, err := s.w.QueryContext(ctx, sql, args...)
	logger(ctx, "QueryContext", sql, ts, args, nil, err)
	return rows, err
}

func (s db) QueryRowContext(ctx context.Context, sql string, args ...any) *sql.Row {
	ts := time.Now()
	row := s.w.QueryRowContext(ctx, sql, args...)
	logger(ctx, "QueryRowContext", sql, ts, args, nil, nil)
	return row
}

// relaceConsecutiveSpaces replaces consecutive spaces with a single space
func relaceConsecutiveSpaces(s string) string {
	for strings.Contains(s, "  ") {
		s = strings.ReplaceAll(s, "  ", " ")
	}
	return s
}

func logger(ctx context.Context, info string, query string, ts time.Time, args any, res sql.Result, err error) {
	if !log.DefaultLogger.Logger.Enabled(ctx, level.Trace) {
		return
	}
	timeSince := time.Since(ts).String()
	query = strings.ReplaceAll(query, "\n", " ")
	query = strings.ReplaceAll(query, "\t", " ")
	query = relaceConsecutiveSpaces(query)
	meta := []any{
		log.String("query", query),
		log.Any("args", args),
		log.String("duration", timeSince),
	}

	if res != nil {
		rows, err := res.RowsAffected()
		if err != nil {
			meta = append(meta, log.String("rows_error", err.Error()))
			log.Error("Rows affected failed", meta...)
		}
		meta = append(meta, log.Int("rows", rows))
	}
	if err != nil {
		meta = append(meta, log.Err(err))
		e := &sqlite.Error{}
		if !errors.As(err, &e) {
			log.Error(fmt.Sprintf("%s failed", info), meta...)
			return
		}
	}
	log.DefaultLogger.Logger.Log(ctx, level.Trace-1, fmt.Sprintf("%s executed", info), meta...)
	log.TraceCtx(ctx, fmt.Sprintf("%s executed", info), meta...)
}
