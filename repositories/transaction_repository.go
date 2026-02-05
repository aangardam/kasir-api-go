package repositories

import (
	"database/sql"
	"fmt"
	"kasir-api/models"
)

type TransactionRepository struct {
	db *sql.DB
}

func NewTransactionRepository(db *sql.DB) *TransactionRepository {
	return &TransactionRepository{db: db}
}

func (repo *TransactionRepository) CreateTransaction(items []models.ChackoutItem) (*models.Transaction, error) {
	tx, err := repo.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	//inisialisasi sub total -> jumlah subtotal keseluruhan
	totalAmount := 0
	//inisialisasi transaction detail -> insert db
	details := make([]models.TransactionDetails, 0)
	//loop items
	for _, item := range items {
		var (
			ProductName string
			ProductID   int
			Price       int
			Stock       int
		)
		//get product id dapatkan price
		err := tx.QueryRow("SELECT id, name, price, stock FROM products where id=$1", item.ProductID).Scan(
			&ProductID,
			&ProductName,
			&Price,
			&Stock,
		)
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("product not found")
		}
		if err != nil {
			return nil, err
		}

		if Stock < item.Quantity {
			return nil, fmt.Errorf("stok produk %s tidak mencukupi (sisa: %d)", ProductName, Stock)
		}
		//hitung current total = quantity * price
		total := item.Quantity * Price
		totalAmount += total

		//kurangi jumlah stock
		_, err = tx.Exec("UPDATE products SET stock = stock - $1 WHERE id = $2", item.Quantity, item.ProductID)
		if err != nil {
			return nil, err
		}

		//item insert ke transaction detail
		details = append(details, models.TransactionDetails{
			ProductID:   ProductID,
			ProductName: ProductName,
			Quantity:    item.Quantity,
			SubTotal:    total,
		})

	}

	//insert ke transaction
	var transactionID int
	err = tx.QueryRow("INSERT INTO transactions (total_amount) VALUES ($1) RETURNING id", totalAmount).Scan(&transactionID)
	if err != nil {
		return nil, err
	}

	query := "INSERT INTO transaction_details (transaction_id, product_id, quantity, subtotal) VALUES "
	var values []interface{}
	counter := 1
	for _, d := range details {
		query += fmt.Sprintf("($%d, $%d, $%d, $%d),", counter, counter+1, counter+2, counter+3)
		values = append(values, transactionID, d.ProductID, d.Quantity, d.SubTotal)
		counter += 4
	}
	// Hapus koma terakhir dan ganti dengan titik koma atau kosong
	query = query[0 : len(query)-1]

	_, err = tx.Exec(query, values...)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &models.Transaction{
		ID:          transactionID,
		TotalAmount: totalAmount,
		Detail:      details,
	}, nil
}
