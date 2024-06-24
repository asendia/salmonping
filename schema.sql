CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS online_listing (
  id uuid NOT NULL DEFAULT uuid_generate_v4(),
  created_at timestamp with time zone DEFAULT now(),
  name text NOT NULL,
  platform text NOT NULL,
  url text NOT NULL,
  PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS schedule (
  id uuid NOT NULL DEFAULT uuid_generate_v4(),
  online_listing_id uuid NOT NULL,
  day_of_week integer NOT NULL,
  opening_time time NOT NULL,
  closing_time time NOT NULL,
  PRIMARY KEY (id),
  FOREIGN KEY (online_listing_id) REFERENCES online_listing(id)
);

CREATE TABLE IF NOT EXISTS salmon_ping (
  id uuid NOT NULL DEFAULT uuid_generate_v4(),
  created_at timestamp with time zone DEFAULT now(),
  status text NOT NULL,
  online_listing_id uuid NOT NULL,
  PRIMARY KEY (id),
  FOREIGN KEY (online_listing_id) REFERENCES online_listing(id)
);

CREATE INDEX IF NOT EXISTS salmon_ping_created_at_idx ON salmon_ping (created_at DESC);
CREATE INDEX IF NOT EXISTS salmon_ping_online_listing_id_idx ON salmon_ping (online_listing_id);
CREATE INDEX IF NOT EXISTS salmon_ping_status_idx ON salmon_ping (status);
CREATE INDEX IF NOT EXISTS salmon_ping_history_filters_idx ON salmon_ping (created_at DESC, status, online_listing_id);

ALTER TABLE online_listing
  ADD COLUMN IF NOT EXISTS enable_ping boolean NOT NULL DEFAULT true;

ALTER TABLE online_listing
  ADD CONSTRAINT url_unique UNIQUE (url);
