package psql

import (
	"fmt"
	"log"
	"order_processing_system/product_service/utils"
	"order_processing_system/user_service/user_utils"

	_ "github.com/jackc/pgx/v5/pgconn"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

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
	GetProductsList() ([]utils.Product, error)
	GetProductByID(id int) (utils.Product, error)
	GetProductQuantity(id int) (utils.Product, error)
	PostProduct(product *utils.Product) error
	PutProduct(product *utils.Product, newProduct utils.Product) error
	DeleteProduct(id int) error
	PostUser(user *user_utils.User) error
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

	_, err = DB.Exec(`CREATE TABLE IF NOT EXISTS product (
		id BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
		name TEXT NOT NULL, description TEXT,
		price NUMERIC(10, 2) NOT NULL,
		stock_quantity INT NOT NULL
		)`)

	if err != nil {
		log.Fatal(err)
	}

	_, err = DB.Exec(`CREATE TABLE IF NOT EXISTS users (
		id BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
		username TEXT NOT NULL,
		email TEXT NOT NULL,
		password_hash TEXT NOT NULL,
		is_admin BOOL NOT NULL,
		created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)`)

	if err != nil {
		log.Fatal(err)
	}

	_, err = DB.Exec(`CREATE TABLE IF NOT EXISTS orders (
		id BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
		user_id BIGINT NOT NULL REFERENCES users (id),
		order_date TIMESTAMPTZ NOT NULL DEFAULT NOW(),
		total_amount NUMERIC(10, 2) NOT NULL
	)`)

	if err != nil {
		log.Fatal(err)
	}

	_, err = DB.Exec(`CREATE TABLE IF NOT EXISTS orders_product (
		order_id BIGINT NOT NULL REFERENCES orders (id),
		product_id BIGINT NOT NULL REFERENCES product (id),
		quantity INT NOT NULL,
		PRIMARY KEY (order_id, product_id)
	)`)

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

func (p *PostgresRepo) GetProductsList() ([]utils.Product, error) {
	var products []utils.Product

	err := p.DB.Select(&products, "SELECT * FROM product")
	if err != nil {
		return nil, err
	}

	return products, nil
}
func (p *PostgresRepo) GetProductByID(id int) (utils.Product, error) {
	var product utils.Product

	err := p.DB.Get(&product, "SELECT * FROM product WHERE id = $1", id)
	if err != nil {
		return utils.Product{}, err
	}

	return product, nil
}

func (p *PostgresRepo) GetProductQuantity(id int) (utils.ProductStock, error) {
	var product utils.ProductStock

	err := p.DB.Get(&product, "SELECT id, stock_quantity FROM product WHERE id = $1", id)
	if err != nil {
		return utils.ProductStock{}, err
	}

	return product, nil
}

func (p *PostgresRepo) PostProduct(product *utils.Product) error {
	err := p.DB.Get(product, "INSERT INTO product (name, description, price, stock_quantity) VALUES ($1, $2, $3, $4) RETURNING *", product.Name, product.Description, product.Price, product.StockQuantity)
	if err != nil {
		return err
	}

	return nil
}

func (p *PostgresRepo) PutProduct(newProduct utils.Product) (utils.Product, error) {
	var updated utils.Product

	err := p.DB.Get(&updated, `
		UPDATE product 
		SET name = $1, description = $2, price = $3, stock_quantity = $4 
		WHERE id = $5 
		RETURNING *`,
		newProduct.Name, newProduct.Description, newProduct.Price, newProduct.StockQuantity, newProduct.ID,
	)

	if err != nil {
		return utils.Product{}, err
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

func (p *PostgresRepo) PostUser(user *user_utils.User) error {
	err := p.DB.Get(user, "INSERT INTO users (username, email, password_hash, is_admin) VALUES ($1, $2, $3, $4) RETURNING *", user.Username, user.Email, user.Password, user.IsAdmin)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}
