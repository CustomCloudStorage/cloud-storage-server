BEGIN;

SET TIME ZONE 'Europe/Moscow';

ALTER TABLE users
    DROP COLUMN name,
    DROP COLUMN email,
    DROP COLUMN password,
    DROP COLUMN role,
    DROP COLUMN storage_limit,
    DROP COLUMN last_update,
    ADD COLUMN created_at TIMESTAMPTZ NOT NULL DEFAULT now();
    
CREATE TABLE profiles (
    user_id INT PRIMARY KEY,
    name TEXT NOT NULL,
    email TEXT NOT NULL UNIQUE,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT fk_profiles_users FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE accounts (
    user_id INT PRIMARY KEY,
    role TEXT NOT NULL,
    storage_limit INT NOT NULL DEFAULT 0,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT fk_accounts_users FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE credentials (
    user_id INT PRIMARY KEY,
    password TEXT NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT fk_credentials_users FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

COMMIT;
