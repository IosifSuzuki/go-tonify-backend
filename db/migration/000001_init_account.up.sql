CREATE TABLE IF NOT EXISTS attachment (
    id SERIAL PRIMARY KEY,
    file_name TEXT NOT NULL,
    path TEXT,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE TABLE IF NOT EXISTS company (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE TABLE IF NOT EXISTS account (
    id SERIAL PRIMARY KEY,
    telegram_id TEXT UNIQUE,
    first_name TEXT NOT NULL,
    middle_name TEXT,
    last_name TEXT NOT NULL,
    nickname TEXT NOT NULL,
    role VARCHAR(32) NOT NULL,
    about_me TEXT NOT NULL,
    gender VARCHAR(127) NOT NULL,
    country TEXT NOT NULL,
    location TEXT NOT NULL,
    avatar_id INT,
    document_id INT,
    company_id INT,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP,
    CONSTRAINT fk_avatar_id FOREIGN KEY(avatar_id) REFERENCES attachment(id) ON DELETE SET NULL,
    CONSTRAINT fk_document_id FOREIGN KEY(document_id) REFERENCES attachment(id) ON DELETE SET NULL,
    CONSTRAINT fk_company_id FOREIGN KEY(company_id) REFERENCES company(id) ON DELETE SET NULL
);
