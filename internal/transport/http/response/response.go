package response

type ErrorResponse struct {
    Errors string `json:"errors"`
}

type AuthResponse struct {
    Token string `json:"token"`
}

type InventoryItem struct {
	Type     string `json:"type"`
	Quantity int    `json:"quantity"`
}

type CoinReceived struct {
	FromUser string `json:"fromUser"`
	Amount   int    `json:"amount"`
}

type CoinSent struct {
	ToUser string `json:"toUser"`
	Amount int    `json:"amount"`
}

type CoinHistory struct {
	Received []CoinReceived `json:"received"`
	Sent     []CoinSent     `json:"sent"`
}

type InfoResponse struct {
	Coins       int            `json:"coins"`
	Inventory   []InventoryItem `json:"inventory"`
	CoinHistory CoinHistory    `json:"coinHistory"`
}
