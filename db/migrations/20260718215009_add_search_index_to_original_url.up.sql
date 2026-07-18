CREATE EXTENSION IF NOT EXISTS pg_trgm;

CREATE INDEX idx_short_urls_original_url_trgm
  ON short_urls USING gin (original_url gin_trgm_ops);
