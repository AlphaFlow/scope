CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE objects
(
    id          UUID PRIMARY KEY,
    db_null_id  UUID,
    num         NUMERIC,
    not_in_json INT
);
