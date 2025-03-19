-- Create "groups" table
CREATE TABLE "groups" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "name" character varying(255) NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT now(),
  "updated_at" timestamptz NOT NULL DEFAULT now(),
  "deleted_at" timestamptz NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_groups_deleted_at" to table: "groups"
CREATE INDEX "idx_groups_deleted_at" ON "groups" ("deleted_at") WHERE (deleted_at IS NOT NULL);
-- Create index "idx_groups_name" to table: "groups"
CREATE INDEX "idx_groups_name" ON "groups" ("name");
-- Create "songs" table
CREATE TABLE "songs" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "group_id" uuid NOT NULL,
  "title" character varying(255) NOT NULL,
  "runtime" integer NOT NULL,
  "lyrics" jsonb NOT NULL,
  "release_date" timestamptz NOT NULL,
  "link" character varying(255) NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT now(),
  "updated_at" timestamptz NOT NULL DEFAULT now(),
  "deleted_at" timestamptz NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_songs_group" FOREIGN KEY ("group_id") REFERENCES "groups" ("id") ON UPDATE NO ACTION ON DELETE CASCADE,
  CONSTRAINT "check_runtime_positive" CHECK (runtime > 0)
);
-- Create index "idx_songs_deleted_at" to table: "songs"
CREATE INDEX "idx_songs_deleted_at" ON "songs" ("deleted_at") WHERE (deleted_at IS NOT NULL);
-- Create index "idx_songs_group_id" to table: "songs"
CREATE INDEX "idx_songs_group_id" ON "songs" ("group_id");
-- Create index "idx_songs_lyrics" to table: "songs"
CREATE INDEX "idx_songs_lyrics" ON "songs" USING gin ("lyrics");
-- Create index "idx_songs_release_date" to table: "songs"
CREATE INDEX "idx_songs_release_date" ON "songs" ("release_date");
-- Create index "idx_songs_title" to table: "songs"
CREATE INDEX "idx_songs_title" ON "songs" ("title");
