CREATE TABLE IF NOT EXISTS "comments" (
  "id" BIGSERIAL PRIMARY KEY,
  "content" TEXT,
  "user_id" BIGINT NOT NULL,
  "post_id" BIGINT NOT NULL,
  "updated_at" TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
  "created_at" TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

ALTER TABLE comments ADD CONSTRAINT fk_comments_user_id FOREIGN KEY (user_id) REFERENCES users (id);
ALTER TABLE comments ADD CONSTRAINT fk_comments_post_id FOREIGN KEY (post_id) REFERENCES posts (id);

CREATE INDEX IF NOT EXISTS idx_comments_user_id ON comments (user_id);
CREATE INDEX IF NOT EXISTS idx_comments_post_id ON comments (post_id);

CREATE OR REPLACE TRIGGER trig_b_u_updated_at_to_now
BEFORE UPDATE ON comments
FOR EACH ROW
EXECUTE PROCEDURE fn_set_updated_at_to_now();