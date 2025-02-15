package merch

import (
	"log/slog"
	"strings"

	"github.com/justcgh9/merch_store/internal/models/inventory"
	"github.com/justcgh9/merch_store/internal/models/transaction"
	"github.com/justcgh9/merch_store/internal/services"
)

type MerchRepo interface {
	BuyStuff(username, item string, cost int) error
	GetInventory(username string) (inventory.Inventory, error)
	GetBalance(username string) (inventory.Balance, error)
	GetHistory(username string) (transaction.TransactionHistory, error)
}

type MerchService struct {
	log       *slog.Logger
	merchRepo MerchRepo
}

func New(log *slog.Logger, merchRepo MerchRepo) *MerchService {
	return &MerchService{
		log:       log,
		merchRepo: merchRepo,
	}
}

func (m *MerchService) Buy(username, item string) error {
	const op = "services.merch.Buy"

	log := m.log.With(
		slog.String("op", op),
		slog.String("username", username),
	)

	log.Info("attempt to buy item", slog.String("item", item))

	cost, ok := inventory.Prices.Get(item)
	if !ok {
		log.Error("item does not exist", slog.String("item", item))
		return services.NonExistingItemError
	}

	item = strings.ReplaceAll(item, "-", "_")

	err := m.merchRepo.BuyStuff(username, item, cost)
	if err != nil {
		log.Error("buy did not succeed", slog.String("err", err.Error()))
		return services.UnsuccessfulBuyError
	}

	log.Info("item bought successfully", slog.String("item", item))

	return nil
}

func (m *MerchService) Informate(username string) (inventory.Info, error) {
	const op = "services.merch.Informate"

	log := m.log.With(
		slog.String("op", op),
		slog.String("username", username),
	)

	log.Info("attempt to get information")

	inv, err := m.merchRepo.GetInventory(username)
	log.Info("inventory", slog.Any("inv", inv), slog.Any("err", err))
	if err != nil {
		log.Error("error accessing inventory", slog.String("err", err.Error()))
		return inventory.Info{}, services.GetInventoryError
	}

	balance, err := m.merchRepo.GetBalance(username)
	if err != nil {
		log.Error("error accessing balance", slog.String("err", err.Error()))
		return inventory.Info{}, services.GetBalanceError
	}

	history, err := m.merchRepo.GetHistory(username)
	if err != nil {
		log.Error("error accessing history", slog.String("err", err.Error()))
		return inventory.Info{}, services.GetHistoryError
	}

	return inventory.Info{
		Inventory:          inv,
		Balance:            balance,
		TransactionHistory: history,
	}, nil
}
