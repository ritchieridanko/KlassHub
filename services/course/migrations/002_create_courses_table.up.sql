CREATE TABLE courses(
  id UUID PRIMARY KEY,
  school_id BIGINT NOT NULL,

  school_course_id VARCHAR, -- school-assigned course id (if any)
  name VARCHAR NOT NULL,
  description TEXT,
  course_picture VARCHAR,

  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  deleted_at TIMESTAMPTZ
);

-- Unique school course for active (not deleted) records
CREATE UNIQUE INDEX idx_courses_unique_school_course ON courses(school_id, school_course_id) WHERE school_course_id IS NOT NULL AND deleted_at IS NULL;

-- Index records by school if active (not deleted)
CREATE INDEX idx_courses_school ON courses(school_id) WHERE deleted_at IS NULL;

-- Index records by name if active (not deleted)
CREATE INDEX idx_courses_name ON courses USING GIN(name gin_trgm_ops) WHERE deleted_at IS NULL;
