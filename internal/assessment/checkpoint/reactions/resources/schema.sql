-- Read-only schema supplied with the reactions exercise.
CREATE TABLE reactions (
  user_id INTEGER NOT NULL,
  post_id INTEGER NOT NULL,
  value   TEXT    NOT NULL CHECK (value IN ('like', 'dislike'))
);
