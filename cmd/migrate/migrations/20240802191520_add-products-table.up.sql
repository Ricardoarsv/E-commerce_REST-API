CREATE TABLE IF NOT EXISTS products (
    ID SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    image VARCHAR(255) NOT NULL,
    quantity INT NOT NULL,
    price FLOAT NOT NULL,
    create_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);