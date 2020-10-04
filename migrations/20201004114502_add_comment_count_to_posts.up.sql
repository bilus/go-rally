ALTER TABLE posts ADD COLUMN comment_count INTEGER NOT NULL DEFAULT 0;

UPDATE posts AS p
SET comment_count = (
    SELECT count(*)
    FROM comments AS c
    WHERE c.post_id = p.id
);
