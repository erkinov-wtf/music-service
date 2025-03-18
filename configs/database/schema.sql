-- Creating the groups table
CREATE TABLE IF NOT EXISTS groups
(
    id           UUID           NOT NULL DEFAULT gen_random_uuid(),
    name         VARCHAR(255)   NOT NULL,
    created_at   TIMESTAMPTZ    NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ    NOT NULL DEFAULT NOW(),
    deleted_at   TIMESTAMPTZ,

    CONSTRAINT groups_pkey PRIMARY KEY (id)
);

CREATE INDEX IF NOT EXISTS idx_groups_name ON groups(name);
CREATE INDEX IF NOT EXISTS idx_groups_deleted_at ON groups(deleted_at) WHERE deleted_at IS NOT NULL;

-- Creating the songs table
CREATE TABLE IF NOT EXISTS songs
(
    id           UUID           NOT NULL DEFAULT gen_random_uuid(),
    group_id     UUID           NOT NULL,
    title        VARCHAR(255)   NOT NULL,
    runtime      INT            NOT NULL,
    lyrics       JSONB          NOT NULL,
    release_date TIMESTAMPTZ    NOT NULL,
    link         VARCHAR(255)   NOT NULL,
    created_at   TIMESTAMPTZ    NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ    NOT NULL DEFAULT NOW(),
    deleted_at   TIMESTAMPTZ,

    CONSTRAINT songs_pkey PRIMARY KEY (id),
    CONSTRAINT fk_songs_group FOREIGN KEY (group_id) REFERENCES groups (id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_songs_group_id ON songs(group_id);
CREATE INDEX IF NOT EXISTS idx_songs_title ON songs(title);
CREATE INDEX IF NOT EXISTS idx_songs_release_date ON songs(release_date);
CREATE INDEX IF NOT EXISTS idx_songs_deleted_at ON songs(deleted_at) WHERE deleted_at IS NOT NULL;

CREATE INDEX IF NOT EXISTS idx_songs_lyrics ON songs USING GIN (lyrics);

ALTER TABLE songs ADD CONSTRAINT check_runtime_positive CHECK (runtime > 0);

-- Adding trigger for updated_at timestamp
CREATE OR REPLACE FUNCTION update_modified_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_groups_modtime
    BEFORE UPDATE ON groups
    FOR EACH ROW
    EXECUTE FUNCTION update_modified_column();

CREATE TRIGGER update_songs_modtime
    BEFORE UPDATE ON songs
    FOR EACH ROW
    EXECUTE FUNCTION update_modified_column();