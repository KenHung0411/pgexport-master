CREATE TABLE rates (
    id BIGSERIAL PRIMARY KEY,
    origin varchar(16) NOT NULL,
    destination varchar(16) NOT NULL,
    container_type integer NOT NULL,
    carrier_scac varchar(16) NOT NULL,
    provider_id bigint NOT NULL,
    price float NOT NULL,
    service_fee float NOT NULL,
    service_code varchar(191),
    effective_date timestamp with time zone NOT NULL,
    expiry_date timestamp with time zone NOT NULL,
    created_at timestamp with time zone NOT NULL DEFAULT now(),
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    version bigint NOT NULL DEFAULT 0,
    contract_number varchar(32),
    commodity varchar(32),
    demurrage text,
    detention text,
    outport text,
    remarks text,
    included text,
    subject_to text,
    port_of_loading varchar(16),
    port_of_discharge varchar(16),
    service_fee_effective_date timestamp with time zone,
    service_fee_expiry_date timestamp with time zone,
    service_mode varchar(16) NOT NULL DEFAULT '',
    promotion_discount float,
    promotion_effective_date timestamp with time zone,
    promotion_expiry_date timestamp with time zone
);

CREATE INDEX index_rates_origin ON rates (origin);
CREATE INDEX index_rates_destination ON rates (destination);
CREATE INDEX index_rates_container_type ON rates (container_type);
CREATE INDEX index_rates_effective_date ON rates (effective_date);
CREATE INDEX index_rates_expiry_date ON rates (expiry_date);
CREATE INDEX index_rates_commodity ON rates ((upper(commodity)));
CREATE INDEX index_rates_port_of_loading ON rates (port_of_loading);
CREATE INDEX index_rates_port_of_discharge ON rates (port_of_discharge);
CREATE INDEX index_rates_service_fee_effective_date ON rates (service_fee_effective_date);
CREATE INDEX index_rates_service_fee_expiry_date ON rates (service_fee_expiry_date);
CREATE INDEX index_rates_service_mode ON rates (service_mode);
CREATE INDEX index_rates_promotion_discount ON rates (promotion_discount);
CREATE INDEX index_rates_promotion_effective_date ON rates (promotion_effective_date);
CREATE INDEX index_rates_promotion_expiry_date ON rates (promotion_expiry_date);
