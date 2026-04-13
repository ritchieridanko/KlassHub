CREATE TABLE classes(
  id UUID PRIMARY KEY,
  course_id UUID NOT NULL,
  school_id BIGINT NOT NULL,

  school_course_id VARCHAR, -- school-assigned course id (if any)
  name VARCHAR NOT NULL,
  description TEXT,
  course_picture VARCHAR,
  syllabus_url VARCHAR,
  min_score SMALLINT,
  academic_year VARCHAR,
  total_students INT NOT NULL DEFAULT 0,
  status VARCHAR NOT NULL DEFAULT 'scheduled', -- e.g. scheduled, ongoing, finished, etc.

  created_by UUID NOT NULL,

  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  deleted_at TIMESTAMPTZ,

  FOREIGN KEY(course_id) REFERENCES courses(id) ON DELETE CASCADE
);

-- Unique school course for active (not deleted) records
CREATE UNIQUE INDEX idx_classes_unique_school_course ON classes(school_id, school_course_id) WHERE school_course_id IS NOT NULL AND deleted_at IS NULL;

-- Index records by course if active (not deleted)
CREATE INDEX idx_classes_course ON classes(course_id) WHERE deleted_at IS NULL;

-- Index records by school if active (not deleted)
CREATE INDEX idx_classes_school ON classes(school_id) WHERE deleted_at IS NULL;

-- Index records by name if active (not deleted)
CREATE INDEX idx_classes_name ON classes USING GIN(name gin_trgm_ops) WHERE deleted_at IS NULL;

-- Index records by status if active (not deleted)
CREATE INDEX idx_classes_status ON classes(status) WHERE deleted_at IS NULL;
