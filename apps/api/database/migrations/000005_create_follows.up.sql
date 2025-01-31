CREATE TABLE IF NOT EXISTS "follows" (
  "following_user_id" BIGINT NOT NULL,
  "followed_user_id" BIGINT NOT NULL,
  "created_at" TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY ("following_user_id", "followed_user_id")
);

ALTER TABLE follows ADD CONSTRAINT fk_follows_following_user_id FOREIGN KEY (following_user_id) REFERENCES users (id) ON DELETE CASCADE;
ALTER TABLE follows ADD CONSTRAINT fk_follows_followed_user_id FOREIGN KEY (followed_user_id) REFERENCES users (id) ON DELETE CASCADE;