CREATE TABLE "urls" (
  -- "id" bigserial PRIMARY KEY,
  "hash_id" varchar(16) NOT NULL UNIQUE,
  "url" varchar NOT NULL
);

CREATE INDEX ON "urls" ("hash_id");
