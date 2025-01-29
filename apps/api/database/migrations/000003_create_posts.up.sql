CREATE TABLE IF NOT EXISTS "posts" (
  "id" BIGSERIAL PRIMARY KEY,
  "title" VARCHAR,
  "content" TEXT,
  "user_id" BIGINT NOT NULL,
  "updated_at" TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
  "created_at" TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

ALTER TABLE posts ADD CONSTRAINT fk_posts_user_id FOREIGN KEY (user_id) REFERENCES users (id);

CREATE INDEX IF NOT EXISTS idx_posts_user_id ON posts (user_id);

CREATE OR REPLACE TRIGGER trig_b_u_updated_at_to_now
BEFORE UPDATE ON posts
FOR EACH ROW
EXECUTE PROCEDURE fn_set_updated_at_to_now();