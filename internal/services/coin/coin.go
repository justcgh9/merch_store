package coin

import (
	"log/slog"

	"github.com/justcgh9/merch_store/internal/services"
)

type CoinRepo interface {
	TransferMoney(to, from string, amount int) error
}

type CoinService struct {
	log      *slog.Logger
	coinRepo CoinRepo
}

func New(log *slog.Logger, coinRepo CoinRepo) *CoinService {
	return &CoinService{
		log:      log,
		coinRepo: coinRepo,
	}
}

func (c *CoinService) Send(from, to string, amount int) error {
	if amount <= 0 {
		return services.TransferZeroMoneyError
	}
	return c.coinRepo.TransferMoney(to, from, amount)
}
