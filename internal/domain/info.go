package domain

type InventoryInfo struct {
	Type     string
	Quantity int
}

type CoinHistoryEntry struct {
	FromUser string
	Amount   int
}

type CoinSentEntry struct {
	ToUser string
	Amount int
}

type CoinHistory struct {
	Received []CoinHistoryEntry
	Sent     []CoinSentEntry
}

type InfoResponse struct {
	Coins       int
	Inventory   []InventoryInfo
	CoinHistory CoinHistory
}
