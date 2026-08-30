CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- id is not a generated key: it IS the user_id auth-api-be issues in the JWT,
-- same pattern as practiq-be's user_profiles. A lookup is a plain PK read,
-- no join needed to resolve "who is this token".
CREATE TABLE IF NOT EXISTS campus_profiles (
    id           VARCHAR(255) PRIMARY KEY,
    profile_type VARCHAR(30) NOT NULL DEFAULT 'student' CHECK (profile_type IN ('teacher', 'student')),
    full_name    VARCHAR(150) DEFAULT '',
    avatar_url   TEXT DEFAULT '',
    bio          TEXT DEFAULT '',
    created_at   TIMESTAMP DEFAULT NOW(),
    updated_at   TIMESTAMP DEFAULT NOW()
);
