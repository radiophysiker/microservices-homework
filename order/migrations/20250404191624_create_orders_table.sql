-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS orders (
    uuid UUID PRIMARY KEY,
    user_uuid UUID NOT NULL,
    total_price DECIMAL(10,2) NOT NULL DEFAULT 0.00,
    transaction_uuid UUID,
    payment_method TEXT,
    status TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);
-- items table for order positions
CREATE TABLE IF NOT EXISTS order_items (
    order_uuid UUID NOT NULL REFERENCES orders(uuid) ON DELETE CASCADE,
    part_uuid UUID NOT NULL,
    quantity INT NOT NULL DEFAULT 1,
    PRIMARY KEY (order_uuid, part_uuid)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS order_items;
DROP TABLE IF EXISTS orders;
-- +goose StatementEnd
