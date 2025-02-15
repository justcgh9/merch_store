package coin_test

import (
	"errors"
	"log/slog"
	"testing"

	"github.com/justcgh9/merch_store/internal/services"
	"github.com/justcgh9/merch_store/internal/services/coin"
	"github.com/justcgh9/merch_store/internal/services/coin/mocks"
	"github.com/stretchr/testify/assert"
)

func TestCoinService_Send(t *testing.T) {
	logger := slog.Default()
	coinRepo := mocks.NewCoinRepo(t)
	service := coin.New(logger, coinRepo)

	t.Run("success", func(t *testing.T) {
		coinRepo.On("TransferMoney", "toUser", "fromUser", 100).Return(nil).Once()

		err := service.Send("fromUser", "toUser", 100)
		assert.NoError(t, err)
		coinRepo.AssertExpectations(t)
	})

	t.Run("error when sending zero money", func(t *testing.T) {
		err := service.Send("fromUser", "toUser", 0)
		assert.ErrorIs(t, err, services.TransferZeroMoneyError)
	})

	t.Run("error when sending negative money", func(t *testing.T) {
		err := service.Send("fromUser", "toUser", -50)
		assert.ErrorIs(t, err, services.TransferZeroMoneyError)
	})

	t.Run("error from coin repo", func(t *testing.T) {
		repoErr := errors.New("transfer error")
		coinRepo.On("TransferMoney", "toUser", "fromUser", 100).Return(repoErr).Once()

		err := service.Send("fromUser", "toUser", 100)
		assert.ErrorIs(t, err, repoErr)
		coinRepo.AssertExpectations(t)
	})
}
