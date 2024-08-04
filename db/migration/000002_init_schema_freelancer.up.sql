CREATE TABLE IF NOT EXISTS freelancer(
    id SERIAL PRIMARY KEY,
    telegram_id TEXT UNIQUE,
    first_name TEXT,
    middle_name TEXT,
    last_name TEXT,
    gender VARCHAR(127),
    country TEXT,
    city TEXT,
    avatar_path TEXT,
    document_path TEXT,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP
)