CREATE TABLE IF NOT EXISTS account_seen
(
    viewer_account_id INT,
    viewed_account_id INT,
    rating INT,
    CONSTRAINT fk_viewer_account_id FOREIGN KEY(viewer_account_id) REFERENCES account(id) ON DELETE SET NULL,
    CONSTRAINT fk_viewed_account_id FOREIGN KEY(viewed_account_id) REFERENCES account(id) ON DELETE SET NULL
)