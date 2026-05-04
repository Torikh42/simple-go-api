package models

// CoinBalance mewakili data saldo user
type CoinBalance struct {
	// Di Go: Nama field diawali kapital agar bisa di-export ke JSON
	// Di JSON: Kita ingin huruf kecil (snake_case atau camelCase)
	Username string `json:"username"`
	Amount   int64  `json:"amount"`
}
