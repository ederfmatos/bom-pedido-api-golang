DROP INDEX IF EXISTS ids_merchants_domain;
DROP INDEX IF EXISTS ids_merchants_tenant_id;
DROP INDEX IF EXISTS ids_merchant_address_merchant_id;
DROP INDEX IF EXISTS ids_merchant_opening_hour_merchant_id;
DROP INDEX IF EXISTS ids_admin_email;
DROP INDEX IF EXISTS ids_customer_email_tenant;
DROP INDEX IF EXISTS ids_product_tenant;
DROP INDEX IF EXISTS ids_order_customer;
DROP INDEX IF EXISTS idx_transactions_order;
DROP INDEX IF EXISTS ids_merchant_payment_gateway_configs_merchant_id;

DROP TABLE IF EXISTS transaction_refunds;
DROP TABLE IF EXISTS pix_transactions;
DROP TABLE IF EXISTS transactions;
DROP TABLE IF EXISTS order_history;
DROP TABLE IF EXISTS order_items;
DROP TABLE IF EXISTS orders;
DROP TABLE IF EXISTS products;
DROP TABLE IF EXISTS customers;
DROP TABLE IF EXISTS admins;
DROP TABLE IF EXISTS merchant_payment_gateway_configs;
DROP TABLE IF EXISTS merchant_opening_hour;
DROP TABLE IF EXISTS merchant_address;
DROP TABLE IF EXISTS merchants;

DROP TYPE IF EXISTS transaction_type;