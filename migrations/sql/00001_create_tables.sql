-- +goose Up
-- SQL in this section is executed when the migration is applied.

CREATE TABLE stats (
    stat_id SERIAL PRIMARY KEY,

    checksum_query varchar NOT NULL,

    occured_at timestamp without time zone NOT NULL
);

CREATE INDEX stats_checksum_query_idx ON stats(checksum_query);

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP TABLE stats;
