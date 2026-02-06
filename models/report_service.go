package models

type DailyReport struct {
	TotalRevenue   int            `json:"total_revenue"`
	TotalTransaksi int            `json:"total_transaksi"`
	ProdukTerlaris ProdutTetlaris `json:"produk_terlaris"`
}

type ProdutTetlaris struct {
	Nama       string `json:"nama"`
	QtyTerjual int    `json:"qty_terjual"`
}

type Report struct {
	Transactions []TransactionReport `json:"transactions"`
}

type TransactionReport struct {
	ID              int     `json:"id"`
	TransactionDate string  `json:"transaction_date"`
	TotalAmount     int     `json:"total_amount"`
	Items           []Items `json:"items"`
}

type Items struct {
	ProductId   int    `json:"product_id"`
	ProductName string `json:"product_name"`
	Category    string `json:"category"`
	Quantity    int    `json:"quantity"`
	Price       int    `json:"price"`
	SubTotal    int    `json:"subtotal"`
}
