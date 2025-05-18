CREATE TABLE IF NOT EXISTS registrations (
    email         VARCHAR(255) PRIMARY KEY,
    password      VARCHAR(255) NOT NULL,
    code          VARCHAR(6)   NOT NULL,
    created_at    TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    last_sent_at  TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);