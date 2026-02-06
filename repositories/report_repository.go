package repositories

import (
	"database/sql"
	"kasir-api/models"
)

type ReportRepository struct {
	db *sql.DB
}

func NewReportRepository(db *sql.DB) *ReportRepository {
	return &ReportRepository{db: db}
}

func (repo *ReportRepository) GetDailyReport() (*models.DailyReport, error) {
	var report models.DailyReport

	// Query 1: Total Revenue & Total Transaksi
	summaryQuery := `
		SELECT 
			COALESCE(SUM(total_amount), 0) as total_revenue, 
			COUNT(id) as total_transaksi 
		FROM transactions 
		WHERE DATE(created_at) = CURRENT_DATE`

	err := repo.db.QueryRow(summaryQuery).Scan(&report.TotalRevenue, &report.TotalTransaksi)
	if err != nil {
		return nil, err
	}

	// Query 2: Produk Terlaris
	bestSellerQuery := `
		SELECT p.name, SUM(td.quantity) as total_qty
		FROM transaction_details td
		JOIN products p ON td.product_id = p.id
		JOIN transactions t ON td.transaction_id = t.id
		WHERE DATE(t.created_at) = CURRENT_DATE
		GROUP BY p.name
		ORDER BY total_qty DESC
		LIMIT 1`

	err = repo.db.QueryRow(bestSellerQuery).Scan(
		&report.ProdukTerlaris.Nama,
		&report.ProdukTerlaris.QtyTerjual,
	)

	if err == sql.ErrNoRows {
		report.ProdukTerlaris.Nama = "Belum ada data"
		report.ProdukTerlaris.QtyTerjual = 0
		return &report, nil
	}

	if err != nil {
		return nil, err
	}

	return &report, nil
}

func (repo *ReportRepository) GetReport(startDate, endDate string) (*models.Report, error) {
	var report models.Report

	// Query 1: Transaksi
	transactionsQuery := `
		SELECT 
			t.id, 
			t.created_at as transaction_date,
			t.total_amount, 
			td.product_id, 
			p.name, 
			td.quantity,  
			td.subtotal,
			c.name as category_name 
		FROM transactions t
		JOIN transaction_details td ON t.id = td.transaction_id
		JOIN products p ON td.product_id = p.id
		JOIN categories c ON p.category_id = c.id
		WHERE DATE(t.created_at) BETWEEN $1 AND $2
		ORDER BY t.id DESC`

	rows, err := repo.db.Query(transactionsQuery, startDate, endDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	transactionsMap := make(map[int]*models.TransactionReport)

	for rows.Next() {
		var (
			transactionID   int
			transactionDate string
			totalAmount     int
			item            models.Items
		)
		// fmt.Println(transactionID, totalAmount)
		err := rows.Scan(
			&transactionID,
			&transactionDate,
			&totalAmount,
			&item.ProductId,
			&item.ProductName,
			&item.Quantity,
			&item.SubTotal,
			&item.Category,
		)
		if err != nil {
			return nil, err
		}
		trx, exists := transactionsMap[transactionID]
		if !exists {
			trx = &models.TransactionReport{
				ID:              transactionID,
				TotalAmount:     totalAmount,
				TransactionDate: transactionDate,
				Items:           []models.Items{},
			}
			transactionsMap[transactionID] = trx
		}

		// append item ke transaksi
		trx.Items = append(trx.Items, item)
	}

	for _, trx := range transactionsMap {
		report.Transactions = append(report.Transactions, *trx)
	}

	return &report, nil
}
