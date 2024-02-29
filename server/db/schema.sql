CREATE TABLE IF NOT EXISTS "schema_migrations" (version varchar(128) primary key);
CREATE TABLE challenges (
  id INTEGER PRIMARY KEY,
  wallet_address CHAR(42) NOT NULL,
  token CHAR(16) NOT NULL UNIQUE,
  expired_at TIMESTAMP NOT NULL
);
CREATE TABLE companies (
  uuid CHAR(36) PRIMARY KEY,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  api_key CHAR(48) NOT NULL,
  admin_account_uuid CHAR(36)
);
CREATE TABLE accounts (
  uuid CHAR(36) PRIMARY KEY,
  company_uuid CHAR(36),
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  wallet_address CHAR(42) NOT NULL UNIQUE,
  metadata JSON
);
-- Dbmate schema migrations
INSERT INTO "schema_migrations" (version) VALUES
  ('20231122185055'),
  ('20240221213521'),
  ('20240229221005');
