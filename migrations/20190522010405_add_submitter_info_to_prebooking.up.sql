ALTER TABLE prebookings ADD COLUMN submit_user_id bigint;
ALTER TABLE prebookings ADD COLUMN source_platform_id bigint;

CREATE INDEX index_prebookings_submit_user_id ON prebookings(submit_user_id);
CREATE INDEX index_prebookings_source_platform_id ON prebookings(source_platform_id);
