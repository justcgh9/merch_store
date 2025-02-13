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


	tx, err := s.conn.Begin()
	if err != nil {
		return fmt.Errorf("%s %v", op, err)
	}
	defer tx.Rollback()

	query := `
	INSERT INTO Users (username, password)
	VALUES ($1, $2);
	`

	_, err = tx.Exec(query, user.Username, user.Password)
	if err != nil {
		return fmt.Errorf("%s %v", op, err)
	}

	query = `
	INSERT INTO Balance (username)
	VALUES ($1);
	`

	_, err = tx.Exec(query, user.Username)
	if err != nil {
		return fmt.Errorf("%s %v", op, err)
	}

	query = `
	INSERT INTO Inventory (username)
	VALUES ($1);
	`

	_, err = tx.Exec(query, user.Username)
	if err != nil {
		return fmt.Errorf("%s %v", op, err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("%s: failed to commit transaction: %v", op, err)
	}

	return nil
}

func (s *Storage) TransferMoney(to, from string, amount int) error {
	const op = "storage.postgres.TransferMoney"

	tx, err := s.conn.Begin()
	if err != nil {
		return fmt.Errorf("%s: begin transaction: %w", op, err)
	}
	defer tx.Rollback()
	
	result, err := tx.Exec(`
        UPDATE balance
        SET balance = balance - $1
        WHERE username = $2 AND balance >= $1
    `, amount, from)
	if err != nil {
		return fmt.Errorf("%s: deduct from sender: %w", op, err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("%s: insufficient funds", op)
	}
	
	result, err = tx.Exec(`
        UPDATE balance
        SET balance = balance + $1
        WHERE username = $2
    `, amount, to)
	if err != nil {
		return fmt.Errorf("%s: add to recipient: %w", op, err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("%s: recipient does not exist", op)
	}
	
	_, err = tx.Exec(`
        INSERT INTO history (from_user, to_user, amount, created_at)
        VALUES ($1, $2, $3, NOW())
    `, from, to, amount)
	if err != nil {
		return fmt.Errorf("%s: insert into history: %w", op, err)
	}
	
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("%s: commit transaction: %w", op, err)
	}

	return nil
}
