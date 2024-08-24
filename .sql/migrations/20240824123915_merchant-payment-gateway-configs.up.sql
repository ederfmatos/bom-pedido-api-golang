CREATE TABLE merchant_payment_gateway_configs
(
    id          SERIAL      NOT NULL PRIMARY KEY,
    merchant_id VARCHAR(36) NOT NULL,
    gateway     VARCHAR(50) NOT NULL,
    credentials TEXT        NOT NULL,
    CONSTRAINT fk_merchant_payment_gateway_configs_merchant FOREIGN KEY (merchant_id) REFERENCES merchants (id)
);
CREATE INDEX ids_merchant_payment_gateway_configs_merchant_id ON merchant_payment_gateway_configs (merchant_id);