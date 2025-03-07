package postgres

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/justcgh9/merch_store/internal/models/inventory"
	"github.com/justcgh9/merch_store/internal/models/transaction"
	"github.com/justcgh9/merch_store/internal/models/user"
	"github.com/justcgh9/merch_store/internal/storage"
	"github.com/lib/pq"
)

type PgxIface interface {
	Begin(context.Context) (pgx.Tx, error)
	Query(context.Context, string, ...interface{}) (pgx.Rows, error)
	QueryRow(context.Context, string, ...interface{}) pgx.Row
	Close()
}

type Storage struct {
	conn    PgxIface
	timeout time.Duration
}

func New(connString string, timeout time.Duration) *Storage {
	const op = "storage.postgres.New"

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	conn, err := pgxpool.New(ctx, connString)
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

func (s *Storage) Close() {
	s.conn.Close()
}

func (s *Storage) GetUser(username string) (user.User, error) {
	const op = "storage.postgres.GetUser"

	ctx, cancel := context.WithTimeout(context.Background(), s.timeout)
	defer cancel()

	var u user.User

	query := `
	SELECT username, password 
	FROM Users 
	WHERE username = $1;
	`

	err := s.conn.QueryRow(ctx, query, username).Scan(&u.Username, &u.Password)
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

	ctx, cancel := context.WithTimeout(context.Background(), s.timeout)
	defer cancel()

	tx, err := s.conn.Begin(ctx)
	if err != nil {
		return fmt.Errorf("%s %v", op, err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	query := `
	INSERT INTO Users (username, password)
	VALUES ($1, $2);
	`

	_, err = tx.Exec(ctx, query, user.Username, user.Password)
	if err != nil {
		return fmt.Errorf("%s %v", op, err)
	}

	query = `
	INSERT INTO Balance (username)
	VALUES ($1);
	`

	_, err = tx.Exec(ctx, query, user.Username)
	if err != nil {
		return fmt.Errorf("%s %v", op, err)
	}

	query = `
	INSERT INTO Inventory (username)
	VALUES ($1);
	`

	_, err = tx.Exec(ctx, query, user.Username)
	if err != nil {
		return fmt.Errorf("%s %v", op, err)
	}

	if err = tx.Commit(ctx); err != nil {
		return fmt.Errorf("%s: failed to commit transaction: %v", op, err)
	}

	return nil
}

func (s *Storage) TransferMoney(to, from string, amount int) error {
	const op = "storage.postgres.TransferMoney"

	ctx, cancel := context.WithTimeout(context.Background(), s.timeout)
	defer cancel()

	tx, err := s.conn.Begin(ctx)
	if err != nil {
		return fmt.Errorf("%s: begin transaction: %w", op, err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	result, err := tx.Exec(ctx, `
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

	result, err = tx.Exec(ctx, `
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

	_, err = tx.Exec(ctx, `
        INSERT INTO history (from_user, to_user, amount, created_at)
        VALUES ($1, $2, $3, NOW())
    `, from, to, amount)
	if err != nil {
		return fmt.Errorf("%s: insert into history: %w", op, err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("%s: commit transaction: %w", op, err)
	}

	return nil
}

func (s *Storage) BuyStuff(username, item string, cost int) error {
	const op = "storage.postgres.BuyStuff"

	ctx, cancel := context.WithTimeout(context.Background(), s.timeout)
	defer cancel()

	tx, err := s.conn.Begin(ctx)
	if err != nil {
		return fmt.Errorf("%s: begin transaction: %w", op, err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	result, err := tx.Exec(ctx, `
        UPDATE balance
        SET balance = balance - $1
        WHERE username = $2 AND balance >= $1
    `, cost, username)
	if err != nil {
		return fmt.Errorf("%s: deduct balance: %w", op, err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("%s: insufficient funds", op)
	}

	query := fmt.Sprintf(`
        UPDATE inventory
        SET %s = %s + 1
        WHERE username = $1
    `, pq.QuoteIdentifier(item), pq.QuoteIdentifier(item))

	result, err = tx.Exec(ctx, query, username)
	if err != nil {
		return fmt.Errorf("%s: update inventory: %w", op, err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("%s: user does not exist in inventory", op)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("%s: commit transaction: %w", op, err)
	}

	return nil
}

func (s *Storage) GetInventory(username string) (inventory.Inventory, error) {
	const op = "storage.postgres.GetInventory"

	ctx, cancel := context.WithTimeout(context.Background(), s.timeout)
	defer cancel()

	row := s.conn.QueryRow(ctx, `
        SELECT t_shirt, cup, book, pen, powerbank, hoody, umbrella, socks, wallet, pink_hoody
        FROM inventory
        WHERE username = $1
    `, username)

	var counts [10]int
	err := row.Scan(
		&counts[0], &counts[1], &counts[2], &counts[3], &counts[4],
		&counts[5], &counts[6], &counts[7], &counts[8], &counts[9],
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, storage.ErrUserDoesNotExist
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	itemNames := []string{
		"t_shirt", "cup", "book", "pen", "powerbank",
		"hoody", "umbrella", "socks", "wallet", "pink_hoody",
	}

	var inv inventory.Inventory
	for i, quantity := range counts {
		if quantity > 0 {
			inv = append(inv, inventory.Item{
				Type:     itemNames[i],
				Quantity: quantity,
			})
		}
	}

	return inv, nil
}

func (s *Storage) GetBalance(username string) (inventory.Balance, error) {
	const op = "storage.postgres.GetBalance"

	ctx, cancel := context.WithTimeout(context.Background(), s.timeout)
	defer cancel()

	var balance inventory.Balance

	err := s.conn.QueryRow(ctx, `SELECT balance FROM balance WHERE username = $1`, username).Scan(&balance)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, storage.ErrUserDoesNotExist
		}
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return balance, nil
}

func (s *Storage) GetHistory(username string) (transaction.TransactionHistory, error) {
	const op = "storage.postgres.GetHistory"

	ctx, cancel := context.WithTimeout(context.Background(), s.timeout)
	defer cancel()

	var history transaction.TransactionHistory

	rows, err := s.conn.Query(ctx, `
		SELECT from_user, to_user, amount
		FROM history
		WHERE from_user = $1 OR to_user = $1
	`, username)
	if err != nil {
		return transaction.TransactionHistory{}, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	for rows.Next() {
		var from, to string
		var amount int

		if err := rows.Scan(&from, &to, &amount); err != nil {
			return transaction.TransactionHistory{}, fmt.Errorf("%s: %w", op, err)
		}

		if to == username {
			history.Recieved = append(history.Recieved, transaction.Recieved{
				From:   from,
				Amount: amount,
			})
		}

		if from == username {
			history.Sent = append(history.Sent, transaction.Sent{
				To:     to,
				Amount: amount,
			})
		}
	}

	if err := rows.Err(); err != nil {
		return transaction.TransactionHistory{}, fmt.Errorf("%s: %w", op, err)
	}

	return history, nil
}
