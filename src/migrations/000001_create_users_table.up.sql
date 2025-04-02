CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    email TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL,
    role TEXT NOT NULL,
    storage_limit INT NOT NULL DEFAULT 0,
    last_update TEXT NOT NULL
);

CREATE INDEX idx_users_email ON users(email);