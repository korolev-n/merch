package e2e

const baseURL = "http://localhost:8080"

type AuthResponse struct {
	Token string `json:"token"`
}

type InventoryItem struct {
	Type     string `json:"type"`
	Quantity int    `json:"quantity"`
}

type InfoResponse struct {
	Coins     int             `json:"coins"`
	Inventory []InventoryItem `json:"inventory"`
}
