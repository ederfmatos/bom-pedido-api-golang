CREATE TABLE IF NOT EXISTS product_categories
(
    id          VARCHAR(36)  NOT NULL PRIMARY KEY,
    name        VARCHAR(255) NOT NULL,
    description TEXT         NOT NULL,
    tenant_id   VARCHAR(30)  NOT NULL,
    created_at  TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT uk_product_category_name UNIQUE (name, tenant_id)
);

CREATE INDEX IF NOT EXISTS ids_product_category_tenant ON product_categories (tenant_id);