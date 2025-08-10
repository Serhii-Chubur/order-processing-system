CREATE TABLE IF NOT EXISTS product (
  id BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
  name VARCHAR(255) NOT NULL,
  description TEXT,
  price NUMERIC(10, 2) NOT NULL,
  stock_quantity INT NOT NULL
);

INSERT INTO product (name, description, price, stock_quantity) VALUES
  ('Laptop Pro X', 'High-performance laptop with 16GB RAM and 512GB SSD', 1200.00, 50),
  ('27-inch 4K Monitor', 'Ultra HD monitor for professional use', 450.00, 75),
  ('USB-C Hub 7-in-1', 'Multi-port USB-C hub with HDMI and card reader', 49.99, 300),
  ('Noise-Cancelling Headphones', 'Premium over-ear headphones with active noise cancellation', 199.99, 100),
  ('Smartwatch Series 5', 'Fitness tracker and smartwatch with heart rate monitor', 250.00, 120),
  ('Portable SSD 1TB', 'External solid-state drive for fast data transfer', 110.00, 90),
  ('Webcam Full HD', '1080p webcam with built-in microphone', 55.00, 180),
  ('Gaming Chair Elite', 'Ergonomic gaming chair with lumbar support', 280.00, 60);
