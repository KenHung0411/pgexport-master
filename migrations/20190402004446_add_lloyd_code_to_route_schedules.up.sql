ALTER TABLE route_schedules ADD COLUMN lloyd_code varchar(128);
CREATE INDEX idx_route_schedules_lloyd_code ON route_schedules(lloyd_code);
