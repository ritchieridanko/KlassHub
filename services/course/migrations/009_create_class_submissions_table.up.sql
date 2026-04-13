CREATE TABLE class_submissions(
  class_assignment_id UUID NOT NULL,
  user_id UUID NOT NULL,
  PRIMARY KEY(class_assignment_id, user_id),

  auth_id BIGINT NOT NULL,
  school_id BIGINT NOT NULL,

  group_id UUID,
  url VARCHAR,
  score SMALLINT,
  attempt INT NOT NULL DEFAULT 0,

  starts_at TIMESTAMPTZ,
  ends_at TIMESTAMPTZ,
  CHECK(starts_at <= ends_at),

  submitted_at TIMESTAMPTZ,
  submitted_by UUID,
  submitted_by_name VARCHAR,

  reviewed_at TIMESTAMPTZ,
  reviewed_by UUID,
  reviewed_by_name VARCHAR,

  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  deleted_at TIMESTAMPTZ,

  FOREIGN KEY(class_assignment_id) REFERENCES class_assignments(id) ON DELETE CASCADE
);

-- Index records by user if active (not deleted)
CREATE INDEX idx_class_submissions_user ON class_submissions(user_id) WHERE deleted_at IS NULL;

-- Index records by auth if active (not deleted)
CREATE INDEX idx_class_submissions_auth ON class_submissions(auth_id) WHERE deleted_at IS NULL;

-- Index records by school if active (not deleted)
CREATE INDEX idx_class_submissions_school ON class_submissions(school_id) WHERE deleted_at IS NULL;

-- Index records by group if exists and active (not deleted)
CREATE INDEX idx_class_submissions_group ON class_submissions(group_id) WHERE group_id IS NOT NULL AND deleted_at IS NULL;
