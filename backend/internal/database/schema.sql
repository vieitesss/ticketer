-- Create stores table
CREATE TABLE IF NOT EXISTS stores (
    id UUID PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE
);

-- Create products table
CREATE TABLE IF NOT EXISTS products (
    id UUID PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    store_id UUID NOT NULL REFERENCES stores(id) ON DELETE RESTRICT,
    UNIQUE(name, store_id)
);

-- Create receipts table
CREATE TABLE IF NOT EXISTS receipts (
    id UUID PRIMARY KEY,
    store_id UUID NOT NULL REFERENCES stores(id) ON DELETE RESTRICT,
    discounts NUMERIC(10, 2) DEFAULT 0,
    receipt_hash VARCHAR(64) UNIQUE,
    bought_date DATE NOT NULL
);

-- Create items table (line items on receipts)
CREATE TABLE IF NOT EXISTS items (
    id UUID PRIMARY KEY,
    receipt_id UUID NOT NULL REFERENCES receipts(id) ON DELETE CASCADE,
    product_id UUID NOT NULL REFERENCES products(id) ON DELETE RESTRICT,
    quantity NUMERIC(10, 3) NOT NULL,
    price_paid NUMERIC(10, 2) NOT NULL
);

-- Create indexes for better query performance
CREATE INDEX IF NOT EXISTS idx_receipts_store_id ON receipts(store_id);
CREATE INDEX IF NOT EXISTS idx_receipts_bought_date ON receipts(bought_date DESC);
CREATE INDEX IF NOT EXISTS idx_receipts_hash ON receipts(receipt_hash);
CREATE INDEX IF NOT EXISTS idx_products_store_id ON products(store_id);
CREATE INDEX IF NOT EXISTS idx_products_name ON products(name);
CREATE INDEX IF NOT EXISTS idx_items_receipt_id ON items(receipt_id);
CREATE INDEX IF NOT EXISTS idx_items_product_id ON items(product_id);
