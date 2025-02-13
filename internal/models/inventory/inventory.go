package inventory

import "github.com/justcgh9/merch_store/internal/models/transaction"

type Info struct {
	Balance            Balance                        `json:"coins"`
	Inventory          Inventory                      `json:"inventory"`
	TransactionHistory transaction.TransactionHistory `json:"coinHistory"`
}

type Inventory = []Item

type Item struct {
	Type     string `json:"type"`
	Quantity int    `json:"quantity"`
}

type Balance = int

type priceMap struct {
	prices map[string]int
}

func (p *priceMap) Get(item string) (int, bool) {
	value, exists := p.prices[item]
	return value, exists
}

var Prices = &priceMap{
	prices: map[string]int{
		"t-shirt":    80,
		"cup":        20,
		"book":       50,
		"pen":        10,
		"powerbank":  200,
		"hoody":      300,
		"umbrella":   200,
		"socks":      10,
		"wallet":     50,
		"pink-hoody": 500,
	},
}
