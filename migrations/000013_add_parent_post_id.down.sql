DROP INDEX IF EXISTS idx_forum_posts_parent_post_id;
ALTER TABLE forum_posts DROP COLUMN IF EXISTS parent_post_id;
