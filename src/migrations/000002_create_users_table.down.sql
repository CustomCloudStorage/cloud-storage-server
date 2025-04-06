BEGIN;

ALTER TABLE users
    ADD COLUMN name TEXT,
    ADD COLUMN email TEXT,
    ADD COLUMN password TEXT,
    ADD COLUMN role TEXT,
    ADD COLUMN storage_limit INT NOT NULL DEFAULT 0,
    ADD COLUMN last_update TEXT,
    DROP COLUMN created_at;

DROP TABLE IF EXISTS profiles;
DROP TABLE IF EXISTS accounts;
DROP TABLE IF EXISTS credentials;

COMMIT;
