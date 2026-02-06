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

func (h *ReportHandler) HandleReport(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.GetReport(w, r)
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

// GetReport - GET /api/report?start_date=YYYY-MM-DD&end_date=YYYY-MM-DD
func (h *ReportHandler) GetReport(w http.ResponseWriter, r *http.Request) {
	startDate := r.URL.Query().Get("start_date")
	endDate := r.URL.Query().Get("end_date")

	if startDate == "" || endDate == "" {
		// http.Error(w, "Missing parameters", http.StatusBadRequest)
		utils.ResponseError(w, http.StatusBadRequest, "Missing parameters")
		return
	}

	if startDate > endDate {
		// http.Error(w, "Invalid parameters", http.StatusBadRequest)
		utils.ResponseError(w, http.StatusBadRequest, "Invalid parameters")
		return
	}

	report, err := h.service.GetReport(startDate, endDate)
	if err != nil {
		utils.ResponseError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	utils.ResponseSuccess(w, report, http.StatusOK, "Report retrieved successfully")
}
