CREATE TABLE IF NOT EXISTS client (
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
    company_id INT,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE TABLE IF NOT EXISTS company (
    id SERIAL PRIMARY KEY,
    name TEXT,
    description TEXT,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP
);

ALTER TABLE client
    ADD CONSTRAINT fk_company_id
    FOREIGN KEY(company_id)
    REFERENCES company(id)
    ON DELETE SET NULL;

