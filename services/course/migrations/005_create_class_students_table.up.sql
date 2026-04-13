CREATE TABLE class_students(
  class_id UUID NOT NULL,
  user_id UUID NOT NULL,
  PRIMARY KEY(class_id, user_id),

  auth_id BIGINT NOT NULL,
  school_id BIGINT NOT NULL,

  school_student_id VARCHAR, -- school-assigned student id (if any)
  name VARCHAR NOT NULL,
  profile_picture VARCHAR,
  final_score SMALLINT,

  assigned_by UUID NOT NULL,

  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  deleted_at TIMESTAMPTZ,

  FOREIGN KEY(class_id) REFERENCES classes(id) ON DELETE CASCADE
);

-- Index records by user if active (not deleted)
CREATE INDEX idx_class_students_user ON class_students(user_id) WHERE deleted_at IS NULL;

-- Index records by auth if active (not deleted)
CREATE INDEX idx_class_students_auth ON class_students(auth_id) WHERE deleted_at IS NULL;

-- Index records by school if active (not deleted)
CREATE INDEX idx_class_students_school ON class_students(school_id) WHERE deleted_at IS NULL;
