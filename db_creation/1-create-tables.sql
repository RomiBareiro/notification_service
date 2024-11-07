-- Create extension pgcrypto if it does not already exist
CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE SCHEMA notification_service;

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TYPE notification_service.notification_type AS ENUM ('STATUS', 'NEWS', 'MARKETING');

CREATE TABLE notification_service.notifications (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    recipient VARCHAR(255) NOT NULL,
    notification_type notification_service.notification_type NOT NULL,
    counter INTEGER,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);


CREATE TABLE notification_service.rate_limit_rules (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    notification_type notification_service.notification_type NOT NULL,
    max_count INTEGER NOT NULL,
    duration NUMERIC NOT NULL
);