CREATE TABLE IF NOT EXISTS prebookings (
    id varchar(64) PRIMARY KEY,
    prebooking JSONB NOT NULL
) WITH (
    OIDS=FALSE
);

CREATE INDEX idx_prebookings_prebooking ON prebookings USING gin (prebooking);

