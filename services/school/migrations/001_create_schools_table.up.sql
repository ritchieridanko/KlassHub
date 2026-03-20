CREATE TABLE schools(
  id BIGSERIAL PRIMARY KEY,

  npsn VARCHAR,
  npsn_verified_at TIMESTAMPTZ,
  CHECK(npsn IS NOT NULL OR npsn_verified_at IS NULL),

  name VARCHAR NOT NULL,
  level VARCHAR NOT NULL, -- e.g. sd, smp, sma, etc.
  ownership VARCHAR NOT NULL, -- e.g. public, private, etc.
  profile_picture VARCHAR,
  profile_banner VARCHAR,
  accreditation VARCHAR,
  established_at TIMESTAMPTZ,

  province VARCHAR NOT NULL,
  city_regency VARCHAR NOT NULL,
  district VARCHAR NOT NULL,
  subdistrict VARCHAR NOT NULL,
  street VARCHAR NOT NULL,
  postcode VARCHAR NOT NULL,

  phone VARCHAR,
  email VARCHAR,
  website VARCHAR,

  timezone VARCHAR NOT NULL,

  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  deleted_at TIMESTAMPTZ
);

-- Unique npsn for verified and active (not deleted) records
CREATE UNIQUE INDEX idx_schools_unique_npsn ON schools(npsn) WHERE npsn_verified_at IS NOT NULL AND deleted_at IS NULL;
