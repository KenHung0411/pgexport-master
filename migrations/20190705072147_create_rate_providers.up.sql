CREATE TABLE rate_providers (
    id BIGSERIAL PRIMARY KEY,
    code varchar(32) NOT NULL,
    name varchar(191) NOT NULL,
    version bigint DEFAULT 0
);

CREATE INDEX index_rate_providers_code ON rate_providers (code);

