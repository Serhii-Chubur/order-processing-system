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
