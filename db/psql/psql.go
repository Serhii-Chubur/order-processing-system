package psql

import (
	"errors"
	"fmt"
	"log"
	"order_processing_system/order_service/order_utils/models"
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
	GetUserByEmail(email string) (user_utils.User, error)
	GetUserByID(id int) (user_utils.User, error)
	GetUserInfo(email string) (user_utils.UserInfo, error)
	PutUser(user *user_utils.UserInput, id int) error
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
		name VARCHAR(255) NOT NULL,
		description TEXT,
		price NUMERIC(10, 2) NOT NULL,
		stock_quantity INT NOT NULL
		)`)

	if err != nil {
		log.Fatal(err)
	}

	_, err = DB.Exec(`CREATE TABLE IF NOT EXISTS users (
		id BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
		username VARCHAR(255) NOT NULL,
		email VARCHAR(255) NOT NULL,
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
		status VARCHAR(255) NOT NULL,
		total_amount NUMERIC(10, 2) NOT NULL
	)`)

	if err != nil {
		log.Fatal(err)
	}

	_, err = DB.Exec(`CREATE TABLE IF NOT EXISTS order_product (
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

func (p *PostgresRepo) DecreaseProductStock(productID int, quantity int) error {
	_, err := p.DB.Exec("UPDATE product SET stock_quantity = stock_quantity - $1 WHERE id = $2", quantity, productID)
	if err != nil {
		return err
	}
	return nil
}

func (p *PostgresRepo) IncreaseProductStock(productID int, quantity int) error {
	_, err := p.DB.Exec("UPDATE product SET stock_quantity = stock_quantity + $1 WHERE id = $2", quantity, productID)
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

func (p *PostgresRepo) GetUserByEmail(email string) (user_utils.User, error) {
	var user user_utils.User
	err := p.DB.Get(&user, "SELECT * FROM users WHERE email = $1", email)
	if err != nil {
		return user_utils.User{}, err
	}
	return user, nil
}

func (p *PostgresRepo) GetUserInfo(email string) (user_utils.UserInfo, error) {
	var user user_utils.UserInfo
	err := p.DB.Get(&user, "SELECT id, username, email, created_at, is_admin FROM users WHERE email = $1", email)
	if err != nil {
		return user_utils.UserInfo{}, err
	}
	return user, nil
}

func (p *PostgresRepo) GetUserById(id int) (user_utils.User, error) {
	var user user_utils.User
	err := p.DB.Get(&user, "SELECT * FROM users WHERE id = $1", id)
	if err != nil {
		return user_utils.User{}, err
	}
	return user, nil
}

func (p *PostgresRepo) PutUser(user *user_utils.UserInput, id int) error {
	_, err := p.DB.Exec("UPDATE users SET username = $1, email = $2, password_hash = $3, is_admin = $4 WHERE id = $5", user.Username, user.Email, user.Password, user.IsAdmin, id)
	if err != nil {
		return err
	}
	return nil
}

func (p *PostgresRepo) PostOrder(order *models.Order) (*models.Order, error) {
	tx, err := p.DB.Beginx()
	if err != nil {
		return nil, err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	err = tx.Get(order, "INSERT INTO orders (user_id, status, total_amount, order_date) VALUES ($1, $2, $3, $4) RETURNING *", order.UserID, order.Status, order.TotalAmount, order.OrderDate)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	for _, product := range order.Products {
		_, err = tx.Exec("INSERT INTO order_product (order_id, product_id, quantity) VALUES ($1, $2, $3)", order.ID, product.ProductID, product.Quantity)
		if err != nil {
			log.Println(err)
			return nil, err
		}
	}
	return order, nil
}

func (p *PostgresRepo) GetOrder(o_id int) (*models.Order, error) {
	var order models.Order
	err := p.DB.Get(&order, "SELECT * FROM orders WHERE id = $1", o_id)
	if err != nil {
		log.Println(err)
		return nil, errors.New("order not found")
	}
	var order_products []models.OrderProduct
	err = p.DB.Select(&order_products, "SELECT product_id, quantity FROM order_product WHERE order_id = $1", o_id)
	if err != nil {
		log.Println(err)
		return nil, errors.New("order not found")
	}
	order.Products = order_products
	return &order, nil
}

func (p *PostgresRepo) GetUserOrders(user_id int) ([]models.Order, error) {
	var orders []models.Order
	err := p.DB.Select(&orders, "SELECT * FROM orders WHERE user_id = $1", user_id)
	if err != nil {
		log.Println(err)
		return nil, errors.New("orders not found")
	}
	return orders, nil
}

func (p *PostgresRepo) PutOrderStatus(o_id int, status string) error {
	res, err := p.DB.Exec("UPDATE orders SET status = $1 WHERE id = $2", status, o_id)
	if err != nil {
		log.Println(err)
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		log.Println(err)
		return err
	}

	if rowsAffected == 0 {
		return errors.New("order not found")
	}
	return nil
}

func (p *PostgresRepo) GetOrderProducts(o_id int) ([]models.OrderProduct, error) {
	var order_products []models.OrderProduct
	err := p.DB.Select(&order_products, "SELECT product_id, quantity FROM order_product WHERE order_id = $1", o_id)
	if err != nil {
		log.Println(err)
		return nil, errors.New("order products not found")
	}
	return order_products, nil
}

// func (p *PostgresRepo) DeleteOrder(o_id int) error {

// 	tx, err := p.DB.Beginx()
// 	if err != nil {
// 		log.Println(err)
// 		return err
// 	}

// 	defer func() {
// 		if err != nil {
// 			tx.Rollback()
// 		} else {
// 			tx.Commit()
// 		}
// 	}()

// 	order, err := p.GetOrder(o_id)
// 	if err != nil {
// 		log.Println(err)
// 		return err
// 	}

// 	for _, product := range order.Products {
// 		_, err = p.DB.Exec("DELETE FROM order_product WHERE order_id = $1 AND product_id = $2", o_id, product.ProductID)
// 		if err != nil {
// 			log.Println(err)
// 			return err
// 		}
// 	}

// 	_, err = p.DB.Exec("DELETE FROM orders WHERE id = $1", o_id)
// 	if err != nil {
// 		log.Println(err)
// 		return err
// 	}
// 	return nil
// }
