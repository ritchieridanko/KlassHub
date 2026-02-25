CREATE TABLE schools(
  id BIGSERIAL PRIMARY KEY,
  auth_id BIGINT NOT NULL,

  npsn VARCHAR,
  npsn_verified_at TIMESTAMPTZ,
  CHECK(npsn IS NOT NULL OR npsn_verified_at IS NULL),

  name VARCHAR NOT NULL,
  school_level VARCHAR NOT NULL, -- e.g. sd, smp, sma, etc.
  school_type VARCHAR NOT NULL, -- e.g. public, private, etc.
  profile_picture VARCHAR,
  profile_banner VARCHAR,
  established_year SMALLINT,
  accreditation_level VARCHAR,

  province VARCHAR NOT NULL,
  city_regency VARCHAR NOT NULL,
  district VARCHAR NOT NULL,
  subdistrict VARCHAR NOT NULL,
  street VARCHAR NOT NULL,
  postcode VARCHAR NOT NULL,

  phone VARCHAR,
  email VARCHAR,
  website VARCHAR,

  is_active BOOLEAN NOT NULL DEFAULT TRUE,
  timezone VARCHAR NOT NULL DEFAULT 'Asia/Jakarta',

  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  deleted_at TIMESTAMPTZ
);

-- Enforce uniqueness of auth id for active (not deleted) records
CREATE UNIQUE INDEX idx_schools_unique_auth ON schools(auth_id) WHERE deleted_at IS NULL;

-- Enforce uniqueness of npsn for active (not deleted) and verified records
CREATE UNIQUE INDEX idx_schools_unique_npsn ON schools(npsn) WHERE npsn_verified_at IS NOT NULL AND deleted_at IS NULL;
