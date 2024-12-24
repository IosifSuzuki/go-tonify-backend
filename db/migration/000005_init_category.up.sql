CREATE TABLE IF NOT EXISTS category (
    id SERIAL PRIMARY KEY,
    title VARCHAR(128) UNIQUE NOT NULL
);

CREATE TABLE IF NOT EXISTS account_category (
    id SERIAL PRIMARY KEY,
    account_id INT NOT NULL,
    category_id INT NOT NULL,
    CONSTRAINT account_category_unique_relationship UNIQUE (account_id, category_id),
    CONSTRAINT fk_account_id FOREIGN KEY (account_id) REFERENCES account(id) ON DELETE CASCADE,
    CONSTRAINT fk_category_id FOREIGN KEY (category_id) REFERENCES category(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS task_category (
    id SERIAL PRIMARY KEY,
    task_id INT NOT NULL,
    category_id INT NOT NULL,
    CONSTRAINT task_category_unique_relationship UNIQUE (task_id, category_id),
    CONSTRAINT fk_category_id FOREIGN KEY (category_id) REFERENCES category(id) ON DELETE CASCADE,
    CONSTRAINT fk_task_id FOREIGN KEY (task_id) REFERENCES task(id) ON DELETE CASCADE
);