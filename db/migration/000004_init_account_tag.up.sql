CREATE TABLE IF NOT EXISTS tag (
    id SERIAL PRIMARY KEY,
    title VARCHAR(128) UNIQUE NOT NULL
);

CREATE TABLE IF NOT EXISTS account_tag (
    id SERIAL PRIMARY KEY,
    account_id INT NOT NULL,
    tag_id INT NOT NULL,
    CONSTRAINT unique_relationship UNIQUE (account_id, tag_id),
    CONSTRAINT fk_account_id FOREIGN KEY (account_id) REFERENCES account(id) ON DELETE CASCADE,
    CONSTRAINT fk_tag_id FOREIGN KEY (tag_id) REFERENCES tag(id) ON DELETE CASCADE
);