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
	var (
		res *models.Transaction
	)

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
		var productName string
		var ProductID int
		var price int
		var stock int
		//get product id dapatkan price
		err := tx.QueryRow("SELECT id, name, price, stock FROM products where id=$1", item.ProductID).Scan(
			&ProductID,
			&productName,
			&price,
			&stock,
		)
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("product not found")
		}
		if err != nil {
			return nil, err
		}
		//hitung current total = quantity * price
		subtotal := item.Quantity * price
		totalAmount += subtotal

		//kurangi jumlah stock
		_, err = tx.Exec("UPDATE products SET stock = stock - $1 WHERE id = $2", item.Quantity, item.ProductID)
		if err != nil {
			return nil, err
		}

		//item insert ke transaction detail
		details = append(details, models.TransactionDetails{
			ProductID:   ProductID,
			ProductName: productName,
			Quantity:    item.Quantity,
			SubTotal:    subtotal,
		})

	}

	//insert ke transaction
	var transactionID int
	err = tx.QueryRow("INSERT INTO transactions (total_amount) VALUES ($1) RETURNING id", totalAmount).Scan(&transactionID)
	if err != nil {
		return nil, err
	}
	fmt.Println(transactionID)

	for i, detail := range details {
		details[i].TransactionID = transactionID
		_, err := tx.Exec("INSERT INTO transaction_details (transaction_id, product_id, quantity, subtotal) VALUES ($1, $2, $3, $4)",
			transactionID,
			detail.ProductID,
			detail.Quantity,
			detail.SubTotal)
		if err != nil {
			return nil, err
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	res = &models.Transaction{
		ID:          transactionID,
		TotalAmount: totalAmount,
		// CreatedAt:   "2022-01-01 00:00:00",
		Detail: details,
	}
	return res, nil
}
