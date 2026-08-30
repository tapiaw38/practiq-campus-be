ALTER TABLE forum_posts
  ADD COLUMN IF NOT EXISTS parent_post_id UUID REFERENCES forum_posts(id) ON DELETE CASCADE;

CREATE INDEX IF NOT EXISTS idx_forum_posts_parent_post_id ON forum_posts(parent_post_id);
