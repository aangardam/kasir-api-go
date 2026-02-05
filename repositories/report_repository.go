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
