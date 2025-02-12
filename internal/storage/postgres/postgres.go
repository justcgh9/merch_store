package postgres

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx"
	"github.com/justcgh9/merch_store/internal/models/user"
	"github.com/justcgh9/merch_store/internal/storage"
)

type Storage struct {
	conn    *pgx.Conn
	timeout time.Duration
}

func New(connString string, timeout time.Duration) *Storage {
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
		conn:    conn,
		timeout: timeout,
	}
}

func (s *Storage) Close() error {
	return s.conn.Close()
}

func (s *Storage) GetUser(username string) (user.User, error) {
	const op = "storage.postgres.GetUser"

	var u user.User

	query := `
	SELECT username, password 
	FROM Users 
	WHERE username = $1;
	`

	err := s.conn.QueryRow(query, username).Scan(&u.Username, &u.Password)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return user.User{}, storage.ErrUserDoesNotExist
		}

		return user.User{}, fmt.Errorf("%s %v", op, err)
	}

	return u, nil
}

func (s *Storage) CreateUser(user user.User) error {
	const op = "storage.postgres.CreateUser"

	query := `
	INSERT INTO Users (username, password)
	VALUES ($1, $2);
	`

	_, err := s.conn.Exec(query, user.Username, user.Password)
	if err != nil {
		return fmt.Errorf("%s %v", op, err)
	}

	return nil
}
