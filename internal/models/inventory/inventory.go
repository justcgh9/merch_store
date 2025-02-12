package inventory

import "github.com/justcgh9/merch_store/internal/models/transaction"

type Info struct {
	Inventory          Inventory                      `json:"inventory"`
	Balance            Balance                        `json:"coins"`
	TransactionHistory transaction.TransactionHistory `json:"coinHistory"`
}

type Inventory []Item

type Item struct {
	Type     string `json:"type"`
	Quantity int    `json:"quantity"`
}

type Balance int
