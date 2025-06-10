package domain

type CoinTransaction struct {
	ID         int
	FromUserID int
	ToUserID   int
	Amount     int
}
