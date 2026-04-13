CREATE TABLE class_sessions(
  id UUID PRIMARY KEY,
  class_id UUID NOT NULL,
  school_id BIGINT NOT NULL,

  session_no SMALLINT NOT NULL,
  topic VARCHAR NOT NULL,
  description TEXT,

  starts_at TIMESTAMPTZ,
  ends_at TIMESTAMPTZ,
  CHECK(starts_at <= ends_at),

  created_by UUID NOT NULL,

  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  deleted_at TIMESTAMPTZ,

  FOREIGN KEY(class_id) REFERENCES classes(id) ON DELETE CASCADE
);

-- Index records by class if active (not deleted)
CREATE INDEX idx_class_sessions_class ON class_sessions(class_id) WHERE deleted_at IS NULL;

-- Index records by school if active (not deleted)
CREATE INDEX idx_class_sessions_school ON class_sessions(school_id) WHERE deleted_at IS NULL;
