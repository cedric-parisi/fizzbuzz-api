-- +goose Up
-- SQL in this section is executed when the migration is applied.

CREATE TABLE stats (
    stat_id SERIAL PRIMARY KEY,

    int1 INTEGER NOT NULL,
    int2 INTEGER NOT NULL,
    max_limit INTEGER NOT NULL,
    str1 varchar NOT NULL,
    str2 varchar NOT NULL,

    occurred_at timestamp without time zone NOT NULL
);

CREATE INDEX stats_all_idx ON stats(int1, int2, max_limit, str1, str2);

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP TABLE stats;
