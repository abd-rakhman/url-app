CREATE TABLE "urls" (
  "hash_id" varchar(16) NOT NULL,
  "url" varchar NOT NULL
);

CREATE INDEX ON "urls" ("hash_id");
