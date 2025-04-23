CREATE TABLE IF NOT EXISTS upload_sessions (
    id UUID PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    folder_id INT REFERENCES folders(id) ON DELETE SET NULL,
    name TEXT NOT NULL,
    extension TEXT NOT NULL,
    total_parts INT NOT NULL,
    total_size BIGINT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS upload_parts (
    session_id UUID NOT NULL REFERENCES upload_sessions(id) ON DELETE CASCADE,
    part_number INT NOT NULL,
    size BIGINT NOT NULL,
    uploaded_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY (session_id, part_number)
);