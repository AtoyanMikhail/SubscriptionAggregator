CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE subscriptions (
    ID PRIMARY KEY DEFAULT uuid_generate_v4(),
    service_name TEXT NOT NULL,
    price INTEGER NOT NULL CHECK (price >= 0),
    user_id UUID NOT NULL,
    start_date DATE NOT NULL,
    end_date DATE
);

CREATE INDEX idx_subscriptions_user_id ON subscriptions(user_id);
CREATE INDEX idx_subscriptions_service_name ON subscriptions(service_name);
CREATE INDEX idx_subscriptions_start_date ON subscriptions(start_date);
CREATE INDEX idx_subscriptions_end_date ON subscriptions(end_date);