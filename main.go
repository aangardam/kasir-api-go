package main

import (
	"encoding/json"
	"fmt"
	"kasir-api/database"
	"kasir-api/handlers"
	"kasir-api/repositories"
	"kasir-api/services"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Port   string `mapstructure:"PORT"`
	DBConn string `mapstructure:"DB_CONN"`
}

func main() {
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if _, err := os.Stat(".env"); err == nil {
		viper.SetConfigFile(".env")
		_ = viper.ReadInConfig()
	}

	config := Config{
		Port:   viper.GetString("PORT"),
		DBConn: viper.GetString("DB_CONN"),
	}

	// Setup database
	db, err := database.InitDB(config.DBConn)
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer db.Close()

	// Products
	productRepo := repositories.NewProductRepository(db)
	productService := services.NewProductService(productRepo)
	productHandler := handlers.NewProductHandler(productService)
	http.HandleFunc("/api/product", productHandler.HandleProducts)
	http.HandleFunc("/api/product/", productHandler.HandleProductByID)

	// Categories
	categoryRepo := repositories.NewCategoryRepository(db)
	categoryService := services.NewCategoryService(categoryRepo)
	categoryHandler := handlers.NewCategoryHandler(categoryService)
	http.HandleFunc("/api/category", categoryHandler.HandleCategorys)
	http.HandleFunc("/api/category/", categoryHandler.HandleCategoryByID)

	// Transaction
	transactionRepo := repositories.NewTransactionRepository(db)
	transactionService := services.NewTransactionService(transactionRepo)
	transactionHandler := handlers.NewTransactionHandler(transactionService)
	http.HandleFunc("/api/checkout", transactionHandler.Checkout)

	// Report
	reportRepo := repositories.NewReportRepository(db)
	reportService := services.NewReportService(reportRepo)
	reportHandler := handlers.NewReportHandler(reportService)
	http.HandleFunc("/api/report/hari-ini", reportHandler.HandleDailyReport)
	http.HandleFunc("/api/report", reportHandler.HandleReport)

	// localhost:8080
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		response := map[string]interface{}{
			"service": "kasir-api-go",
			"version": "v1",
			"endpoints": map[string]interface{}{
				"product": []map[string]string{
					{"method": "GET", "path": "/api/product", "desc": "List products"},
					{"method": "GET", "path": "/api/product/{id}", "desc": "Detail product"},
					{"method": "POST", "path": "/api/product", "desc": "Create product"},
					{"method": "PUT", "path": "/api/product/{id}", "desc": "Update product"},
					{"method": "DELETE", "path": "/api/product/{id}", "desc": "Delete product"},
				},
				"category": []map[string]string{
					{"method": "GET", "path": "/api/category", "desc": "List categories"},
					{"method": "GET", "path": "/api/category/{id}", "desc": "Detail category"},
					{"method": "POST", "path": "/api/category", "desc": "Create category"},
					{"method": "PUT", "path": "/api/category/{id}", "desc": "Update category"},
					{"method": "DELETE", "path": "/api/category/{id}", "desc": "Delete category"},
				},
				"transaction": []map[string]string{
					{"method": "POST", "path": "/api/checkout", "desc": "Checkout"},
				},
				"report": []map[string]string{
					{"method": "GET", "path": "/api/report/hari-ini", "desc": "Daily report"},
					{"method": "GET", "path": "/api/report?start_date=YYYY-MM-DD&end_date=YYYY-MM-DD", "desc": "Report by date range"},
				},
			},
		}

		json.NewEncoder(w).Encode(response)
	})
	fmt.Println("Server running di localhost:" + config.Port)

	err = http.ListenAndServe(":"+config.Port, nil)
	if err != nil {
		fmt.Println("gagal running server")
	}
}
