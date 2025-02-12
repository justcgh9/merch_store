package postgres

import (
	"context"
	"log"

	"github.com/jackc/pgx"
)

type Storage struct {
	conn *pgx.Conn
}

func New(connString string) *Storage {
	const op = "storage.postgres.New"

	config, err := pgx.ParseURI(connString)
	if err != nil {
		log.Fatalf("%s %v", op, err)
	}

	conn, err := pgx.Connect(config)
	if err != nil {
		log.Fatalf("%s %v", op, err)
	}

	err = conn.Ping(context.Background())
	if err != nil {
		log.Fatalf("%s %v", op, err)
	}

	return &Storage{
		conn: conn,
	}
}

func (s *Storage) Close() error {
	return s.conn.Close()
}
