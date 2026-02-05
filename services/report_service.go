package services

import (
	"kasir-api/models"
	"kasir-api/repositories"
)

type ReportService struct {
	// Pastikan nama repository sesuai dengan yang didefinisikan di folder repositories
	repo *repositories.ReportRepository
}

func NewReportService(repo *repositories.ReportRepository) *ReportService {
	return &ReportService{repo: repo}
}

func (s *ReportService) GetDailyReport() (*models.DailyReport, error) {
	// Memanggil fungsi di repository
	return s.repo.GetDailyReport()
}
