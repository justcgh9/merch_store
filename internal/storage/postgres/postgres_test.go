package postgres_test

import (
	"errors"
	"fmt"
	"reflect"
	"testing"
	"time"
	"unsafe"

	"github.com/jackc/pgx/v5"
	"github.com/justcgh9/merch_store/internal/models/user"
	"github.com/justcgh9/merch_store/internal/storage"
	"github.com/justcgh9/merch_store/internal/storage/postgres"
	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/assert"
)

func TestGetUser_Success(t *testing.T) {
	mockConn, err := pgxmock.NewPool()
	assert.NoError(t, err)

	defer mockConn.Close()

	rows := pgxmock.NewRows([]string{"username", "password"}).
		AddRow("testuser", "testpassword")
	mockConn.ExpectQuery("SELECT username, password FROM Users WHERE username = \\$1;").
		WithArgs("testuser").
		WillReturnRows(rows)

	store := &postgres.Storage{}

	setFieldValue(store, "conn", mockConn)
	setFieldValue(store, "timeout", 3*time.Second)

	u, err := store.GetUser("testuser")
	assert.NoError(t, err)
	assert.Equal(t, "testuser", u.Username)
	assert.Equal(t, "testpassword", u.Password)
	assert.NoError(t, mockConn.ExpectationsWereMet())
}

func TestGetUser_NotFound(t *testing.T) {
	mockConn, err := pgxmock.NewPool()
	assert.NoError(t, err)

	defer mockConn.Close()

	store := &postgres.Storage{}

	setFieldValue(store, "conn", mockConn)
	setFieldValue(store, "timeout", 3*time.Second)

	mockConn.ExpectQuery("SELECT username, password FROM Users WHERE username = \\$1;").
		WithArgs("nonexistent").
		WillReturnError(pgx.ErrNoRows)

	_, err = store.GetUser("nonexistent")
	assert.ErrorIs(t, err, storage.ErrUserDoesNotExist)
	assert.NoError(t, mockConn.ExpectationsWereMet())
}

func TestGetUser_QueryError(t *testing.T) {
	mockConn, err := pgxmock.NewPool()
	assert.NoError(t, err)

	defer mockConn.Close()

	store := &postgres.Storage{}

	setFieldValue(store, "conn", mockConn)
	setFieldValue(store, "timeout", 3*time.Second)

	mockConn.ExpectQuery("SELECT username, password FROM Users WHERE username = \\$1;").
		WithArgs("erroruser").
		WillReturnError(errors.New("query error"))

	_, err = store.GetUser("erroruser")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "storage.postgres.GetUser")
	assert.NoError(t, mockConn.ExpectationsWereMet())
}

func TestCreateUser_Success(t *testing.T) {
	mockConn, err := pgxmock.NewPool()
	assert.NoError(t, err)

	defer mockConn.Close()

	store := &postgres.Storage{}

	setFieldValue(store, "conn", mockConn)
	setFieldValue(store, "timeout", 3*time.Second)

	mockConn.ExpectBegin()

	mockConn.ExpectExec("INSERT INTO Users").
		WithArgs("testuser", "hashedpassword").
		WillReturnResult(pgxmock.NewResult("INSERT", 1))

	mockConn.ExpectExec("INSERT INTO Balance").
		WithArgs("testuser").
		WillReturnResult(pgxmock.NewResult("INSERT", 1))

	mockConn.ExpectExec("INSERT INTO Inventory").
		WithArgs("testuser").
		WillReturnResult(pgxmock.NewResult("INSERT", 1))

	mockConn.ExpectCommit()

	err = store.CreateUser(user.User{Username: "testuser", Password: "hashedpassword"})
	assert.NoError(t, err)
	assert.NoError(t, mockConn.ExpectationsWereMet())
}

func TestCreateUser_InsertUserError(t *testing.T) {
	mockConn, err := pgxmock.NewPool()
	assert.NoError(t, err)

	defer mockConn.Close()

	store := &postgres.Storage{}

	setFieldValue(store, "conn", mockConn)
	setFieldValue(store, "timeout", 3*time.Second)

	mockConn.ExpectBegin()
	mockConn.ExpectExec("INSERT INTO Users").
		WithArgs("testuser", "hashedpassword").
		WillReturnError(errors.New("insert error"))
	mockConn.ExpectRollback()

	err = store.CreateUser(user.User{Username: "testuser", Password: "hashedpassword"})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "storage.postgres.CreateUser")
	assert.NoError(t, mockConn.ExpectationsWereMet())
}

func TestTransferMoney_Success(t *testing.T) {
	mockConn, err := pgxmock.NewPool()
	assert.NoError(t, err)

	defer mockConn.Close()

	store := &postgres.Storage{}

	setFieldValue(store, "conn", mockConn)
	setFieldValue(store, "timeout", 3*time.Second)

	mockConn.ExpectBegin()
	mockConn.ExpectExec("UPDATE balance").
		WithArgs(50, "sender").
		WillReturnResult(pgxmock.NewResult("UPDATE", 1))
	mockConn.ExpectExec("UPDATE balance").
		WithArgs(50, "recipient").
		WillReturnResult(pgxmock.NewResult("UPDATE", 1))
	mockConn.ExpectExec("INSERT INTO history").
		WithArgs("sender", "recipient", 50).
		WillReturnResult(pgxmock.NewResult("INSERT", 1))
	mockConn.ExpectCommit()

	err = store.TransferMoney("recipient", "sender", 50)
	assert.NoError(t, err)
	assert.NoError(t, mockConn.ExpectationsWereMet())
}

func TestTransferMoney_InsufficientFunds(t *testing.T) {
	mockConn, err := pgxmock.NewPool()
	assert.NoError(t, err)

	defer mockConn.Close()

	store := &postgres.Storage{}

	setFieldValue(store, "conn", mockConn)
	setFieldValue(store, "timeout", 3*time.Second)

	mockConn.ExpectBegin()
	mockConn.ExpectExec("UPDATE balance").
		WithArgs(50, "sender").
		WillReturnResult(pgxmock.NewResult("UPDATE", 0))
	mockConn.ExpectRollback()

	err = store.TransferMoney("recipient", "sender", 50)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "insufficient funds")
	assert.NoError(t, mockConn.ExpectationsWereMet())
}

func TestTransferMoney_RecipientNotExist(t *testing.T) {
	mockConn, err := pgxmock.NewPool()
	assert.NoError(t, err)

	defer mockConn.Close()

	store := &postgres.Storage{}

	setFieldValue(store, "conn", mockConn)
	setFieldValue(store, "timeout", 3*time.Second)

	mockConn.ExpectBegin()
	mockConn.ExpectExec("UPDATE balance").
		WithArgs(50, "sender").
		WillReturnResult(pgxmock.NewResult("UPDATE", 1))
	mockConn.ExpectExec("UPDATE balance").
		WithArgs(50, "recipient").
		WillReturnResult(pgxmock.NewResult("UPDATE", 0))
	mockConn.ExpectRollback()

	err = store.TransferMoney("recipient", "sender", 50)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "recipient does not exist")
	assert.NoError(t, mockConn.ExpectationsWereMet())
}

func TestBuyStuff_Success(t *testing.T) {
	mockConn, err := pgxmock.NewPool()
	assert.NoError(t, err)

	defer mockConn.Close()

	store := &postgres.Storage{}

	setFieldValue(store, "conn", mockConn)
	setFieldValue(store, "timeout", 3*time.Second)

	mockConn.ExpectBegin()
	mockConn.ExpectExec("UPDATE balance").
		WithArgs(80, "user1").
		WillReturnResult(pgxmock.NewResult("UPDATE", 1))

	mockConn.ExpectExec(`UPDATE inventory SET "t_shirt" = "t_shirt" \+ 1 WHERE username = \$1`).
		WithArgs("user1").
		WillReturnResult(pgxmock.NewResult("UPDATE", 1))
	mockConn.ExpectCommit()

	err = store.BuyStuff("user1", "t_shirt", 80)
	assert.NoError(t, err)
	assert.NoError(t, mockConn.ExpectationsWereMet())
}

func TestGetInventory_Success(t *testing.T) {
	mockConn, err := pgxmock.NewPool()
	assert.NoError(t, err)

	defer mockConn.Close()

	store := &postgres.Storage{}

	setFieldValue(store, "conn", mockConn)
	setFieldValue(store, "timeout", 3*time.Second)

	rows := pgxmock.NewRows([]string{"t_shirt", "cup", "book", "pen", "powerbank", "hoody", "umbrella", "socks", "wallet", "pink_hoody"}).
		AddRow(2, 0, 0, 0, 0, 0, 0, 0, 0, 0)
	mockConn.ExpectQuery("SELECT t_shirt, cup, book, pen, powerbank, hoody, umbrella, socks, wallet, pink_hoody FROM inventory WHERE username = \\$1").
		WithArgs("user1").
		WillReturnRows(rows)

	inv, err := store.GetInventory("user1")
	assert.NoError(t, err)
	assert.Equal(t, 1, len(inv))
	assert.Equal(t, "t_shirt", inv[0].Type)
	assert.Equal(t, 2, inv[0].Quantity)
	assert.NoError(t, mockConn.ExpectationsWereMet())
}

func TestGetBalance_Success(t *testing.T) {
	mockConn, err := pgxmock.NewPool()
	assert.NoError(t, err)

	defer mockConn.Close()

	store := &postgres.Storage{}

	setFieldValue(store, "conn", mockConn)
	setFieldValue(store, "timeout", 3*time.Second)

	rows := pgxmock.NewRows([]string{"balance"}).AddRow(500)
	mockConn.ExpectQuery("SELECT balance FROM balance WHERE username = \\$1").
		WithArgs("user1").
		WillReturnRows(rows)

	bal, err := store.GetBalance("user1")
	assert.NoError(t, err)
	assert.Equal(t, 500, bal)
	assert.NoError(t, mockConn.ExpectationsWereMet())
}

func TestGetHistory_Success(t *testing.T) {
	mockConn, err := pgxmock.NewPool()
	assert.NoError(t, err)

	defer mockConn.Close()

	store := &postgres.Storage{}

	setFieldValue(store, "conn", mockConn)
	setFieldValue(store, "timeout", 3*time.Second)

	rows := pgxmock.NewRows([]string{"from_user", "to_user", "amount"}).
		AddRow("sender1", "user1", 50).
		AddRow("user1", "recipient1", 30)
	mockConn.ExpectQuery("SELECT from_user, to_user, amount FROM history WHERE (.+)").
		WithArgs("user1").
		WillReturnRows(rows)

	hist, err := store.GetHistory("user1")
	assert.NoError(t, err)
	assert.Equal(t, 1, len(hist.Recieved))
}

func setFieldValue(target any, fieldName string, value any) {
	rv := reflect.ValueOf(target)
	for rv.Kind() == reflect.Ptr && !rv.IsNil() {
		rv = rv.Elem()
	}
	if !rv.CanAddr() {
		panic("target must be addressable")
	}
	if rv.Kind() != reflect.Struct {
		panic(fmt.Sprintf(
			"unable to set the '%s' field value of the type %T, target must be a struct",
			fieldName,
			target,
		))
	}
	rf := rv.FieldByName(fieldName)

	reflect.NewAt(rf.Type(), unsafe.Pointer(rf.UnsafeAddr())).Elem().Set(reflect.ValueOf(value))
}
