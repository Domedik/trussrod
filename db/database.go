package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/Domedik/trussrod/settings"
	_ "github.com/lib/pq"
)

type DB struct {
	Conn *sql.DB
}

func getDSN(c *settings.DatabaseConfig) string {
	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		c.Host, c.User, c.Password, c.Name, c.Port, c.SSLMode,
	)
}

func New(c *settings.DatabaseConfig) (*DB, error) {
	dsn := getDSN(c)
	var err error
	conn, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	conn.SetMaxOpenConns(50)
	conn.SetMaxIdleConns(10)
	conn.SetConnMaxLifetime(10 * time.Minute)
	conn.SetConnMaxIdleTime(5 * time.Minute)

	return &DB{
		Conn: conn,
	}, nil
}

func (db *DB) Close() {
	db.Conn.Close()
}

func (db *DB) Prepare(ctx context.Context, query string) (*sql.Stmt, error) {
	return db.Conn.PrepareContext(ctx, query)
}

func (db *DB) Ping(ctx context.Context) error {
	return db.Conn.PingContext(ctx)
}
