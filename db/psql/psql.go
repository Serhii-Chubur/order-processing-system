package psql

import (
	"fmt"
	"log"

	_ "github.com/jackc/pgx/v5/pgconn"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

type Product struct {
	ID            int     `db:"id" json:"id"`
	Name          string  `db:"name" json:"name"`
	Description   string  `db:"description" json:"description"`
	Price         float64 `db:"price" json:"price"`
	StockQuantity int     `db:"stock_quantity" json:"stock"`
}

func (p *Product) Validate() error {
	if p.Name == "" {
		return fmt.Errorf("name is required")
	}
	if p.Price <= 0 {
		return fmt.Errorf("price must be greater than 0")
	}
	if p.StockQuantity <= 0 {
		return fmt.Errorf("stock quantity must be greater than 0")
	}
	return nil
}

type ProductStock struct {
	ID            int `db:"id" json:"id"`
	StockQuantity int `db:"stock_quantity" json:"stock"`
}

type PSQLConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

type PostgresRepo struct {
	DB *sqlx.DB
}

type PostgresCRUD interface {
	GetProductsList() ([]Product, error)
	GetProductByID(id int) (Product, error)
	GetProductQuantity(id int) (Product, error)
	PostProduct(product *Product) error
	PutProduct(product *Product, newProduct Product) error
	DeleteProduct(id int) error
}

func ConnectPSQL(config PSQLConfig) *sqlx.DB {
	connString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", config.User, config.Password, config.Host, config.Port, config.DBName)
	DB, err := sqlx.Connect("pgx", connString)
	if err != nil {
		log.Fatal(err)
	}

	if err = DB.Ping(); err != nil {
		log.Fatal(err)
	}

	_, err = DB.Exec(`CREATE TABLE IF NOT EXISTS product (id BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY, name TEXT NOT NULL, description TEXT, price NUMERIC(10, 2) NOT NULL, stock_quantity INT NOT NULL)`)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to PSQL")
	return DB
}

func NewPSQLRepo(db *sqlx.DB) *PostgresRepo {
	return &PostgresRepo{
		DB: db,
	}
}

func (p *PostgresRepo) GetProductsList() ([]Product, error) {
	var products []Product

	err := p.DB.Select(&products, "SELECT * FROM product")
	if err != nil {
		return nil, err
	}

	return products, nil
}
func (p *PostgresRepo) GetProductByID(id int) (Product, error) {
	var product Product

	err := p.DB.Get(&product, "SELECT * FROM product WHERE id = $1", id)
	if err != nil {
		return Product{}, err
	}

	return product, nil
}

func (p *PostgresRepo) GetProductQuantity(id int) (ProductStock, error) {
	var product ProductStock

	err := p.DB.Get(&product, "SELECT id, stock_quantity FROM product WHERE id = $1", id)
	if err != nil {
		return ProductStock{}, err
	}

	return product, nil
}

func (p *PostgresRepo) PostProduct(product *Product) error {
	err := p.DB.Get(product, "INSERT INTO product (name, description, price, stock_quantity) VALUES ($1, $2, $3, $4) RETURNING *", product.Name, product.Description, product.Price, product.StockQuantity)
	if err != nil {
		return err
	}

	return nil
}

func (p *PostgresRepo) PutProduct(newProduct Product) (Product, error) {
	var updated Product

	err := p.DB.Get(&updated, `
		UPDATE product 
		SET name = $1, description = $2, price = $3, stock_quantity = $4 
		WHERE id = $5 
		RETURNING *`,
		newProduct.Name, newProduct.Description, newProduct.Price, newProduct.StockQuantity, newProduct.ID,
	)

	if err != nil {
		return Product{}, err
	}

	return updated, nil
}

func (p *PostgresRepo) DeleteProduct(id int) error {
	_, err := p.DB.Exec("DELETE FROM product WHERE id = $1", id)
	if err != nil {
		return err
	}
	return nil
}
