package handlers

import (
	"kasir-api/services"
	"kasir-api/utils"
	"net/http"
)

type ReportHandler struct {
	service *services.ReportService
}

func NewReportHandler(service *services.ReportService) *ReportHandler {
	return &ReportHandler{service: service}
}

func (h *ReportHandler) HandleDailyReport(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.GetDailyReport(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// GetDailyReport - GET /api/report
func (h *ReportHandler) GetDailyReport(w http.ResponseWriter, r *http.Request) {
	report, err := h.service.GetDailyReport()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	utils.ResponseSuccess(w, report, http.StatusOK, "Daily report retrieved successfully")
}
