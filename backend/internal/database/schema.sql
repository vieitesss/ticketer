-- Create receipts table
CREATE TABLE IF NOT EXISTS receipts (
    id UUID PRIMARY KEY,
    store_name VARCHAR(255) NOT NULL,
    discounts NUMERIC(10, 2) DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Create items table
CREATE TABLE IF NOT EXISTS items (
    id UUID PRIMARY KEY,
    receipt_id UUID NOT NULL REFERENCES receipts(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    quantity NUMERIC(10, 3) NOT NULL,
    price NUMERIC(10, 2) NOT NULL
);

-- Create indexes for better query performance
CREATE INDEX IF NOT EXISTS idx_receipts_created_at ON receipts(created_at DESC);
CREATE INDEX IF NOT EXISTS idx_items_receipt_id ON items(receipt_id);
