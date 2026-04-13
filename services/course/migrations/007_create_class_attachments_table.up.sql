CREATE TABLE class_attachments(
  id UUID PRIMARY KEY,
  class_session_id UUID NOT NULL,
  class_id UUID NOT NULL,
  school_id BIGINT NOT NULL,

  title VARCHAR NOT NULL,
  type VARCHAR NOT NULL, -- e.g. document, image, audio, video, url, etc.
  url VARCHAR NOT NULL,

  created_by UUID NOT NULL,

  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  deleted_at TIMESTAMPTZ,

  FOREIGN KEY(class_session_id) REFERENCES class_sessions(id) ON DELETE CASCADE,
  FOREIGN KEY(class_id) REFERENCES classes(id) ON DELETE CASCADE
);

-- Index records by class session if active (not deleted)
CREATE INDEX idx_class_attachments_class_session ON class_attachments(class_session_id) WHERE deleted_at IS NULL;

-- Index records by class if active (not deleted)
CREATE INDEX idx_class_attachments_class ON class_attachments(class_id) WHERE deleted_at IS NULL;

-- Index records by school if active (not deleted)
CREATE INDEX idx_class_attachments_school ON class_attachments(school_id) WHERE deleted_at IS NULL;
