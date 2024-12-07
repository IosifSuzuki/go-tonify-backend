CREATE TABLE IF NOT EXISTS attachment (
    id SERIAL PRIMARY KEY,
    file_name TEXT,
    path TEXT,
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

CREATE TABLE IF NOT EXISTS account (
    id SERIAL PRIMARY KEY,
    telegram_id TEXT UNIQUE,
    first_name TEXT,
    middle_name TEXT,
    last_name TEXT,
    nickname TEXT,
    role VARCHAR(32),
    about_me TEXT,
    gender VARCHAR(127),
    country TEXT,
    location TEXT,
    avatar_id INT,
    document_id INT,
    company_id INT,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP,
    CONSTRAINT fk_avatar_id FOREIGN KEY(avatar_id) REFERENCES attachment(id) ON DELETE SET NULL,
    CONSTRAINT fk_document_id FOREIGN KEY(document_id) REFERENCES attachment(id) ON DELETE SET NULL,
    CONSTRAINT fk_company_id FOREIGN KEY(company_id) REFERENCES company(id) ON DELETE SET NULL
);
