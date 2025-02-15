package merch_test

import (
	"errors"
	"testing"

	"log/slog"

	"github.com/justcgh9/merch_store/internal/models/inventory"
	"github.com/justcgh9/merch_store/internal/models/transaction"
	"github.com/justcgh9/merch_store/internal/services"
	"github.com/justcgh9/merch_store/internal/services/merch"
	"github.com/justcgh9/merch_store/internal/services/merch/mocks"
	"github.com/stretchr/testify/assert"
)

func TestMerchService_Buy(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		username      string
		item          string
		mockBehaviour func(repo *mocks.MerchRepo)
		expectError   error
	}{
		{
			name:     "successful purchase",
			username: "user1",
			item:     "t-shirt",
			mockBehaviour: func(repo *mocks.MerchRepo) {
				repo.On("BuyStuff", "user1", "t_shirt", 80).Return(nil)
			},
			expectError: nil,
		},
		{
			name:          "non-existing item",
			username:      "user1",
			item:          "non-existing-item",
			mockBehaviour: func(repo *mocks.MerchRepo) {},
			expectError:   services.NonExistingItemError,
		},
		{
			name:     "unsuccessful purchase",
			username: "user1",
			item:     "t-shirt",
			mockBehaviour: func(repo *mocks.MerchRepo) {
				repo.On("BuyStuff", "user1", "t_shirt", 80).Return(errors.New("some error"))
			},
			expectError: services.UnsuccessfulBuyError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := mocks.NewMerchRepo(t)
			log := slog.Default()
			service := merch.New(log, repo)

			tt.mockBehaviour(repo)

			err := service.Buy(tt.username, tt.item)

			assert.ErrorIs(t, err, tt.expectError)
		})
	}
}

func TestMerchService_Informate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		username      string
		mockBehaviour func(repo *mocks.MerchRepo)
		expectResult  inventory.Info
		expectError   error
	}{
		{
			name:     "successful informate",
			username: "user1",
			mockBehaviour: func(repo *mocks.MerchRepo) {
				repo.On("GetInventory", "user1").Return(inventory.Inventory{{Type: "t-shirt", Quantity: 2}}, nil)
				repo.On("GetBalance", "user1").Return(100, nil)
				repo.On("GetHistory", "user1").Return(transaction.TransactionHistory{
					Recieved: []transaction.Recieved{{From: "user2", Amount: 50}},
					Sent:     []transaction.Sent{{To: "user3", Amount: 30}},
				}, nil)
			},
			expectResult: inventory.Info{
				Inventory: inventory.Inventory{{Type: "t-shirt", Quantity: 2}},
				Balance:   100,
				TransactionHistory: transaction.TransactionHistory{
					Recieved: []transaction.Recieved{{From: "user2", Amount: 50}},
					Sent:     []transaction.Sent{{To: "user3", Amount: 30}},
				},
			},
			expectError: nil,
		},
		{
			name:     "get inventory error",
			username: "user1",
			mockBehaviour: func(repo *mocks.MerchRepo) {
				repo.On("GetInventory", "user1").Return(nil, errors.New("inventory error"))
			},
			expectResult: inventory.Info{},
			expectError:  services.GetInventoryError,
		},
		{
			name:     "get balance error",
			username: "user1",
			mockBehaviour: func(repo *mocks.MerchRepo) {
				repo.On("GetInventory", "user1").Return(inventory.Inventory{{Type: "t-shirt", Quantity: 2}}, nil)
				repo.On("GetBalance", "user1").Return(0, errors.New("balance error"))
			},
			expectResult: inventory.Info{},
			expectError:  services.GetBalanceError,
		},
		{
			name:     "get history error",
			username: "user1",
			mockBehaviour: func(repo *mocks.MerchRepo) {
				repo.On("GetInventory", "user1").Return(inventory.Inventory{{Type: "t-shirt", Quantity: 2}}, nil)
				repo.On("GetBalance", "user1").Return(100, nil)
				repo.On("GetHistory", "user1").Return(transaction.TransactionHistory{}, errors.New("history error"))
			},
			expectResult: inventory.Info{},
			expectError:  services.GetHistoryError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := mocks.NewMerchRepo(t)
			log := slog.Default()
			service := merch.New(log, repo)

			tt.mockBehaviour(repo)

			result, err := service.Informate(tt.username)

			assert.ErrorIs(t, err, tt.expectError)
			assert.Equal(t, tt.expectResult, result)
		})
	}
}
