DROP INDEX IF EXISTS index_prebookings_submit_user_id;
DROP INDEX IF EXISTS index_prebookings_source_platform_id;

ALTER TABLE prebookings DROP COLUMN IF EXISTS submit_user_id;
ALTER TABLE prebookings DROP COLUMN IF EXISTS source_platform_id;

