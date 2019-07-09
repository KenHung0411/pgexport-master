-- booking summary table
CREATE TABLE booking_summary (
    id bigserial PRIMARY KEY,
    region varchar(128),
    platform_id bigint not null,
    application_number varchar(128) not null,
    carrier varchar(16),
    status integer not null,
    is_deleted boolean,
    created_by bigint,
    updated_by bigint,
    created_at bigint,
    updated_at bigint,
    data jsonb
);
CREATE INDEX idx_booking_summary_platform_id ON booking_summary(platform_id);
CREATE UNIQUE INDEX idx_booking_summary_application_number ON booking_summary USING btree (application_number);
CREATE INDEX idx_booking_summary_carrier ON booking_summary(carrier);
CREATE INDEX idx_booking_summary_status ON booking_summary(status);
CREATE INDEX idx_booking_summary_is_deleted ON booking_summary(is_deleted);
CREATE INDEX idx_booking_summary_created_by ON booking_summary(created_by);
CREATE INDEX idx_booking_summary_created_at ON booking_summary(created_at);
CREATE INDEX idx_booking_summary_updated_by ON booking_summary(updated_by);
CREATE INDEX idx_booking_summary_updated_at ON booking_summary(updated_at);
CREATE INDEX idx_booking_summary_data ON booking_summary USING GIN(data);

-- booking confirm summary table
CREATE TABLE booking_confirm_summary (
    id bigserial PRIMARY KEY,
    region varchar(128),
    application_number varchar(128) not null,
    booking_number varchar(128) not null,
    status integer not null,
    ts bigint,
    data jsonb
);
CREATE UNIQUE INDEX idx_booking_confirm_summary_application_booking_number ON booking_confirm_summary USING btree (application_number, booking_number);
CREATE INDEX idx_booking_confirm_summary_status ON booking_confirm_summary(status);
CREATE INDEX idx_booking_confirm_summary_data ON booking_confirm_summary USING GIN(data);

-- si summary table
CREATE TABLE si_summary (
    id bigserial PRIMARY KEY,
    region varchar(128),
    platform_id bigint not null,
    application_number varchar(128) not null,
    carrier varchar(16),
    status integer not null,
    is_deleted boolean,
    created_by bigint,
    updated_by bigint,
    created_at bigint,
    updated_at bigint,
    data jsonb,
    confirm jsonb
);
CREATE INDEX idx_si_summary_platform_id ON si_summary(platform_id);
CREATE UNIQUE INDEX idx_si_summary_application_number ON si_summary USING btree (application_number);
CREATE INDEX idx_si_summary_carrier ON si_summary(carrier);
CREATE INDEX idx_si_summary_status ON si_summary(status);
CREATE INDEX idx_si_summary_is_deleted ON si_summary(is_deleted);
CREATE INDEX idx_si_summary_created_by ON si_summary(created_by);
CREATE INDEX idx_si_summary_updated_by ON si_summary(updated_by);
CREATE INDEX idx_si_summary_created_at ON si_summary(created_at);
CREATE INDEX idx_si_summary_updated_at ON si_summary(updated_at);
CREATE INDEX idx_si_summary_data ON si_summary USING GIN(data);
CREATE INDEX idx_si_summary_confirm ON si_summary USING GIN(confirm);

-- vgm summary table
CREATE TABLE vgm_summary (
    id bigserial PRIMARY KEY,
    region varchar(128),
    platform_id bigint not null,
    application_number varchar(128) not null,
    carrier varchar(16),
    status integer not null,
    is_deleted boolean,
    created_by bigint,
    updated_by bigint,
    created_at bigint,
    updated_at bigint,
    data jsonb,
    confirm jsonb
);
CREATE INDEX idx_vgm_summary_platform_id ON vgm_summary(platform_id);
CREATE UNIQUE INDEX idx_vgm_summary_application_number ON vgm_summary USING btree (application_number);
CREATE INDEX idx_vgm_summary_carrier ON vgm_summary(carrier);
CREATE INDEX idx_vgm_summary_status ON vgm_summary(status);
CREATE INDEX idx_vgm_summary_is_deleted ON vgm_summary(is_deleted);
CREATE INDEX idx_vgm_summary_created_by ON vgm_summary(created_by);
CREATE INDEX idx_vgm_summary_updated_by ON vgm_summary(updated_by);
CREATE INDEX idx_vgm_summary_created_at ON vgm_summary(created_at);
CREATE INDEX idx_vgm_summary_updated_at ON vgm_summary(updated_at);
CREATE INDEX idx_vgm_summary_data ON vgm_summary USING GIN(data);
CREATE INDEX idx_vgm_summary_confirm ON vgm_summary USING GIN(confirm);

