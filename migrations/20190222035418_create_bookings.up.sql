CREATE TABLE booking (
    id bigserial PRIMARY KEY,
    region varchar(128),
    application_number varchar(128),
    version bigint,
    is_draft boolean,
    platform_id bigint,
    created_by bigint,
    updated_by bigint,
    created_at bigint,
    updated_at bigint,
    data jsonb
);
CREATE UNIQUE INDEX idx_booking_app_version ON booking USING btree (application_number, version);
CREATE INDEX idx_booking_platform_id ON booking(platform_id);
CREATE INDEX idx_booking_created_by ON booking(created_by);
CREATE INDEX idx_booking_updated_by ON booking(updated_by);
CREATE INDEX idx_booking_data ON booking USING GIN(data);

CREATE TABLE booking_confirm (
    id bigserial PRIMARY KEY,
    region varchar(128),
    application_number varchar(128),
    version bigint,
    created_at bigint,
    data jsonb
);
CREATE INDEX idx_booking_confirm_application_number ON booking_confirm(application_number);
CREATE INDEX idx_booking_confirm_data ON booking_confirm USING GIN(data);

CREATE TABLE si (
    id bigserial PRIMARY KEY,
    region varchar(128),
    application_number varchar(128),
    version bigint,
    is_draft boolean,
    platform_id bigint,
    created_by bigint,
    updated_by bigint,
    created_at bigint,
    updated_at bigint,
    data jsonb
);
CREATE UNIQUE INDEX idx_si_app_version ON si USING btree (application_number, version);
CREATE INDEX idx_si_platform_id ON si(platform_id);
CREATE INDEX idx_si_created_by ON si(created_by);
CREATE INDEX idx_si_updated_by ON si(updated_by);
CREATE INDEX idx_si_data ON si USING GIN(data);

CREATE TABLE si_confirm (
    id bigserial PRIMARY KEY,
    region varchar(128),
    application_number varchar(128),
    version bigint,
    created_at bigint,
    data jsonb
);
CREATE INDEX idx_si_confirm_application_number ON si_confirm(application_number);
CREATE INDEX idx_si_confirm_data ON si_confirm USING GIN(data);

CREATE TABLE vgm (
    id bigserial PRIMARY KEY,
    region varchar(128),
    application_number varchar(128),
    version bigint,
    is_draft boolean,
    platform_id bigint,
    created_by bigint,
    updated_by bigint,
    created_at bigint,
    updated_at bigint,
    data jsonb
);
CREATE UNIQUE INDEX idx_vgm_app_version ON vgm USING btree (application_number, version);
CREATE INDEX idx_vgm_platform_id ON vgm(platform_id);
CREATE INDEX idx_vgm_created_by ON vgm(created_by);
CREATE INDEX idx_vgm_updated_by ON vgm(updated_by);
CREATE INDEX idx_vgm_data ON vgm USING GIN(data);

CREATE TABLE vgm_confirm (
    id bigserial PRIMARY KEY,
    region varchar(128),
    application_number varchar(128),
    version bigint,
    created_at bigint,
    data jsonb
);
CREATE INDEX idx_vgm_confirm_application_number ON vgm_confirm(application_number);
CREATE INDEX idx_vgm_confirm_data ON vgm_confirm USING GIN(data);

