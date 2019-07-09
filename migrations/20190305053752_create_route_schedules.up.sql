-- replicate_progress table
CREATE TABLE replicate_progress (
    task varchar(256) NOT NULL,
    table_name varchar(256) NOT NULL,
    progress bigint NOT NULL DEFAULT 0,
    updated_at bigint NOT NULL DEFAULT 0
);
CREATE UNIQUE INDEX idx_replicate_progress_task_table_name
    ON replicate_progress
    USING btree (task, table_name);

-- route_schedules table
CREATE TABLE route_schedules (
    id bigserial PRIMARY KEY,
    carrier varchar(128),
    vessel varchar(128),
    voyage varchar(128),
    last_etd bigint,
    data jsonb
);
CREATE UNIQUE INDEX idx_route_schedules_carrier_vessel_voyage
    ON route_schedules
    USING btree (upper((carrier)::text), upper((vessel)::text), upper((voyage)::text));
CREATE INDEX idx_route_schedules_carrier
    ON route_schedules(carrier);
CREATE INDEX idx_route_schedules_vessel
    ON route_schedules(vessel);
CREATE INDEX idx_route_schedules_voyage
    ON route_schedules(voyage);
CREATE INDEX idx_route_schedules_last_etd
    ON route_schedules(last_etd);
CREATE INDEX idx_route_schedules_data
    ON route_schedules
    USING GIN(data);
