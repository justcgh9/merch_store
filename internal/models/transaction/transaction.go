package transaction

type Recieved struct {
	From   string `json:"fromUser"`
	Amount int    `json:"amount"`
}

type Sent struct {
	To     string `json:"toUser"`
	Amount int    `json:"amount"`
}

type TransactionHistory struct {
	Recieved []Recieved `json:"recieved"`
	Sent     []Sent     `json:"sent"`
}
