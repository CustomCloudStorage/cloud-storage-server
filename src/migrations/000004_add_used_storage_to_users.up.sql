ALTER TABLE accounts
    ADD COLUMN used_storage BIGINT NOT NULL DEFAULT 0;