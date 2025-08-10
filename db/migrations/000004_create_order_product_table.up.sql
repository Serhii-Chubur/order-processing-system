CREATE TABLE IF NOT EXISTS order_product (
	order_id BIGINT NOT NULL REFERENCES orders (id),
	product_id BIGINT NOT NULL REFERENCES product (id),
	quantity INT NOT NULL,
	PRIMARY KEY (order_id, product_id)
);

INSERT INTO order_product (order_id, product_id, quantity) VALUES
	(1, 1, 1),
	(2, 2, 1);