CREATE TABLE IF NOT EXISTS like_account (
    id SERIAL PRIMARY KEY,
    liker_id INT NOT NULL,
    liked_id INT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_liker FOREIGN KEY (liker_id) REFERENCES account(id) ON DELETE CASCADE,
    CONSTRAINT fk_liked FOREIGN KEY (liked_id) REFERENCES account(id) ON DELETE CASCADE,
    CONSTRAINT unique_like UNIQUE (liker_id, liked_id),
    CONSTRAINT no_self_like CHECK (liker_id != liked_id)
);

CREATE TABLE  IF NOT EXISTS dislike_account (
    id SERIAL PRIMARY KEY,
    disliker_id INT NOT NULL,
    disliked_id INT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_disliker FOREIGN KEY (disliker_id) REFERENCES account(id) ON DELETE CASCADE,
    CONSTRAINT fk_disliked FOREIGN KEY (disliked_id) REFERENCES account(id) ON DELETE CASCADE,
    CONSTRAINT unique_dislike UNIQUE (disliker_id, disliked_id),
    CONSTRAINT no_self_dislike CHECK (disliker_id != disliked_id)
);