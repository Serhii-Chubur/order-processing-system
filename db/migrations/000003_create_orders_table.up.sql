CREATE TABLE IF NOT EXISTS orders (
	id BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
	user_id BIGINT NOT NULL REFERENCES users (id),
	order_date TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	status VARCHAR(255) NOT NULL,
	total_amount NUMERIC(10, 2) NOT NULL
);

INSERT INTO orders (user_id, status, total_amount) VALUES
	(1, 'created', 1200.00),
	(2, 'created', 450.00);
