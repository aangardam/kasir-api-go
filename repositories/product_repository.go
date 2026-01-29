package repositories

import (
	"database/sql"
	"errors"
	"kasir-api/models"
)

type ProductRepository struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

func (repo *ProductRepository) GetAll() ([]models.Product, error) {
	// query := "SELECT id, name, price, stock FROM products"
	query := `
		SELECT 
			p.id, 
			p.name, 
			p.price, 
			p.stock, 
			p.category_id,
			c.name AS category_name,
			c.description
		FROM products p
		LEFT JOIN categories c ON c.id = p.category_id`

	rows, err := repo.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	products := make([]models.Product, 0)
	for rows.Next() {
		var (
			p            models.Product
			cID          sql.NullInt64
			cName        sql.NullString
			cDescription sql.NullString
		)
		// err := rows.Scan(&p.ID, &p.Name, &p.Price, &p.Stock)
		err := rows.Scan(
			&p.ID,
			&p.Name,
			&p.Price,
			&p.Stock,
			&cID,
			&cName,
			&cDescription,
		)
		if err != nil {
			return nil, err
		}
		if cID.Valid {
			p.Category = &models.Category{
				ID:          int(cID.Int64),
				Name:        cName.String,
				Description: cDescription.String,
			}
		}
		products = append(products, p)
	}

	return products, nil
}

func (repo *ProductRepository) Create(product *models.Product) error {
	query := "INSERT INTO products (name, price, stock, category_id) VALUES ($1, $2, $3, $4) RETURNING id"
	err := repo.db.QueryRow(query, product.Name, product.Price, product.Stock, product.CategoryID).Scan(&product.ID)
	return err
}

// GetByID - ambil produk by ID
func (repo *ProductRepository) GetByID(id int) (*models.Product, error) {
	// query := "SELECT id, name, price, stock FROM products WHERE id = $1"
	query := `
		SELECT 
			p.id, 
			p.name, 
			p.price, 
			p.stock, 
			p.category_id,
			c.name AS category_name,
			c.description
		FROM products p
		LEFT JOIN categories c ON c.id = p.category_id
		WHERE p.id = $1
	`

	// var p models.Product
	var (
		p            models.Product
		cID          sql.NullInt64
		cName        sql.NullString
		cDescription sql.NullString
	)
	// err := repo.db.QueryRow(query, id).Scan(&p.ID, &p.Name, &p.Price, &p.Stock)
	err := repo.db.QueryRow(query, id).Scan(
		&p.ID,
		&p.Name,
		&p.Price,
		&p.Stock,
		&cID,
		&cName,
		&cDescription,
	)
	if err == sql.ErrNoRows {
		return nil, errors.New("Product Not Found")
	}
	if err != nil {
		return nil, err
	}
	if cID.Valid {
		p.Category = &models.Category{
			ID:          int(cID.Int64),
			Name:        cName.String,
			Description: cDescription.String,
		}
	}

	return &p, nil
}

func (repo *ProductRepository) Update(product *models.Product) error {
	query := "UPDATE products SET name = $1, price = $2, stock = $3, category_id = $4 WHERE id = $5"
	result, err := repo.db.Exec(query, product.Name, product.Price, product.Stock, product.CategoryID, product.ID)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return errors.New("Product Not Found")
	}

	return nil
}

func (repo *ProductRepository) Delete(id int) error {
	query := "DELETE FROM products WHERE id = $1"
	result, err := repo.db.Exec(query, id)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return errors.New("Product Not Found")
	}

	return err
}
