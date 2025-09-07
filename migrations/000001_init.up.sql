CREATE TABLE IF NOT EXISTS orders (
    order_uid TEXT NOT NULL PRIMARY KEY,
    track_number TEXT NOT NULL,
    entry TEXT,
    locale TEXT,
    internal_signature TEXT,
    customer_id TEXT,
    delivery_service TEXT,
    shardKey TEXT,
    sm_id BIGINT,
    date_created TIMESTAMPTZ NOT NULL,
    oof_shard TEXT
);

CREATE TABLE IF NOT EXISTS delivery (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    order_uid TEXT REFERENCES orders(order_uid) ON DELETE CASCADE,
    name TEXT NOT NULL,
    phone TEXT NOT NULL,
    zip TEXT,
    city TEXT NOT NULL,
    address TEXT NOT NULL,
    region TEXT NOT NULL,
    email TEXT
);

CREATE TABLE IF NOT EXISTS payments (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    order_uid TEXT REFERENCES orders(order_uid) ON DELETE CASCADE,
    transaction TEXT NOT NULL,
    request_id TEXT,
    currency TEXT NOT NULL,
    provider TEXT NOT NULL,
    amount INT NOT NULL,
    payment_dt BIGINT,
    bank TEXT,
    delivery_cost INT,
    goods_total INT,
    custom_fee INT
);

CREATE TABLE IF NOT EXISTS items (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    order_uid TEXT REFERENCES orders(order_uid) ON DELETE CASCADE,
    chrt_id BIGINT NOT NULL,
    track_number TEXT,
    price INT,
    rid TEXT,
    name TEXT NOT NULL,
    sale INT NOT NULL,
    size TEXT,
    total_price INT NOT NULL,
    nm_id BIGINT,
    brand TEXT,
    status INT
);