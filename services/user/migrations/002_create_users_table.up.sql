CREATE TABLE users(
  id UUID PRIMARY KEY,
  auth_id BIGINT NOT NULL,
  school_id BIGINT NOT NULL,

  school_user_id VARCHAR,
  role VARCHAR NOT NULL, -- e.g. student, instructor, staff, administrator, etc.
  name VARCHAR NOT NULL,
  nickname VARCHAR,
  birthplace VARCHAR NOT NULL,
  birthdate DATE NOT NULL,
  sex VARCHAR NOT NULL,
  phone VARCHAR,
  profile_picture VARCHAR,
  profile_banner VARCHAR,

  created_by UUID,
  created_by_name VARCHAR,

  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  deleted_at TIMESTAMPTZ,

  FOREIGN KEY(created_by) REFERENCES users(id) ON DELETE CASCADE
);

-- Enforce uniqueness of auth id for active (not deleted) records
CREATE UNIQUE INDEX idx_users_unique_auth ON users(auth_id) WHERE deleted_at IS NULL;

-- Enforce uniqueness of school user id for active (not deleted) records
CREATE UNIQUE INDEX idx_users_unique_school_user ON users(school_id, school_user_id) WHERE school_user_id IS NOT NULL AND deleted_at IS NULL;

-- Index records by school id if active (not deleted)
CREATE INDEX idx_users_school ON users(school_id) WHERE deleted_at IS NULL;

-- Index records by school id and role if active (not deleted)
CREATE INDEX idx_users_school_role ON users(school_id, role) WHERE deleted_at IS NULL;

-- Index records by name if active (not deleted)
CREATE INDEX idx_users_name ON users USING GIN(name gin_trgm_ops) WHERE deleted_at IS NULL;

-- Index records by nickname if exists and active (not deleted)
CREATE INDEX idx_users_nickname ON users USING GIN(nickname gin_trgm_ops) WHERE nickname IS NOT NULL AND deleted_at IS NULL;
