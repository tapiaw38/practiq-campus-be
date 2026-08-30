CREATE TABLE IF NOT EXISTS course_materials (
    id           UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    course_id    UUID NOT NULL REFERENCES courses(id) ON DELETE CASCADE,
    section_id   UUID REFERENCES course_sections(id) ON DELETE SET NULL,
    uploader_id  VARCHAR(255) NOT NULL REFERENCES campus_profiles(id) ON DELETE CASCADE,
    title        VARCHAR(200) NOT NULL,
    description  TEXT NOT NULL DEFAULT '',
    kind         VARCHAR(20) NOT NULL CHECK (kind IN ('file', 'link')),
    -- For kind='file' this is the canonical (private) bucket URL, presigned at
    -- read time; for kind='link' it is the external URL the teacher pasted.
    url          TEXT NOT NULL,
    created_at   TIMESTAMP DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_course_materials_course_id ON course_materials(course_id);
