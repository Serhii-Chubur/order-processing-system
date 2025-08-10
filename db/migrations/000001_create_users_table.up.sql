CREATE TABLE IF NOT EXISTS users (
	id BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
	username VARCHAR(255) NOT NULL,
	email VARCHAR(255) NOT NULL,
	password_hash TEXT NOT NULL,
	is_admin BOOL NOT NULL,
	created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

INSERT INTO users (username, email, password_hash, is_admin) VALUES
	('Admin', 'admin@admin.com', '$2a$10$OEFij9SAtm8JI/7CCCUHyeWfIF4Sc4VuePC9DF/Ou5wBQMNzERuW.', true),
	('Test', 'test@test.com', '$2a$10$kzJnenSnvY8c70qyMUgAqe.n7evaM6yfyuSN9PMMtRsIDe7.CGn..', false);
