-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS orders (
    uuid UUID PRIMARY KEY,
    user_uuid UUID NOT NULL,
    part_uuids UUID[] NOT NULL DEFAULT '{}',
    total_price DECIMAL(10,2) NOT NULL DEFAULT 0.00,
    transaction_uuid UUID,
    payment_method VARCHAR(50),
    status VARCHAR(50) NOT NULL DEFAULT 'PENDING_PAYMENT',
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER IF NOT EXISTS update_orders_updated_at 
    BEFORE UPDATE ON orders 
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TRIGGER IF EXISTS update_orders_updated_at ON orders;
DROP FUNCTION IF EXISTS update_updated_at_column();
DROP TABLE IF EXISTS orders;
-- +goose StatementEnd
