DROP TABLE IF EXISTS users;

CREATE TABLE users
(
    id SERIAL PRIMARY KEY,
    username VARCHAR(255) NOT NULL UNIQUE,
    full_name VARCHAR(255) NOT NULL,
    server_id INTEGER NOT NULL,
    password VARCHAR(255) NOT NULL,
    valid_until TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (server_id) REFERENCES servers (id),
    UNIQUE (username, server_id)
);