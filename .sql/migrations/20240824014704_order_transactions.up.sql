CREATE TABLE transactions
(
    id         VARCHAR(36) PRIMARY KEY NOT NULL,
    order_id   VARCHAR(36)             NOT NULL,
    amount     NUMERIC(6, 2)           NOT NULL,
    status     VARCHAR(20)             NOT NULL,
    created_at TIMESTAMP               NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_transactions_order_id FOREIGN KEY (order_id) REFERENCES orders (id)
);

CREATE INDEX idx_transactions_order ON transactions (order_id);

CREATE TABLE pix_transactions
(
    qr_code         TEXT        NOT NULL,
    qr_code_link    TEXT        NOT NULL,
    payment_gateway VARCHAR(50) NOT NULL
) inherits (transactions);