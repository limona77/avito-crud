package model

type UserInfo struct {
	Coins       int              `json:"coins"`
	Inventory   []*UserInventory `json:"inventory"`
	CoinHistory CoinHistory      `json:"coinHistory"`
}

type UserInventory struct {
	Type   string `json:"type"`
	Amount int    `json:"amount"`
}

type CoinHistory struct {
	Received []*Received `json:"received"`
	Sent     []*Sent     `json:"sent"`
}

type Received struct {
	FromUser string `json:"fromUser" db:"sender"`
	Amount   int    `json:"amount"`
}

type Sent struct {
	ToUser string `json:"toUser" db:"receiver"` // db тег для привязки имени столбца
	Amount int    `json:"amount"`
}
