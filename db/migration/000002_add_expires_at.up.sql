ALTER TABLE urls ADD COLUMN expires_at TIMESTAMP NOT NULL DEFAULT (now() + interval '1 day') ;