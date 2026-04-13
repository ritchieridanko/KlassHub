CREATE TABLE class_assignments(
  id UUID PRIMARY KEY,
  class_session_id UUID NOT NULL,
  class_id UUID NOT NULL,
  school_id BIGINT NOT NULL,

  title VARCHAR NOT NULL,
  type VARCHAR NOT NULL, -- e.g. personal, group, etc.
  url VARCHAR NOT NULL,

  starts_at TIMESTAMPTZ,
  ends_at TIMESTAMPTZ,
  CHECK(starts_at <= ends_at),

  created_by UUID NOT NULL,

  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  deleted_at TIMESTAMPTZ,

  FOREIGN KEY(class_session_id) REFERENCES class_sessions(id) ON DELETE CASCADE,
  FOREIGN KEY(class_id) REFERENCES classes(id) ON DELETE CASCADE
);

-- Index records by class session if active (not deleted)
CREATE INDEX idx_class_assignments_class_session ON class_assignments(class_session_id) WHERE deleted_at IS NULL;

-- Index records by class if active (not deleted)
CREATE INDEX idx_class_assignments_class ON class_assignments(class_id) WHERE deleted_at IS NULL;

-- Index records by school if active (not deleted)
CREATE INDEX idx_class_assignments_school ON class_assignments(school_id) WHERE deleted_at IS NULL;
