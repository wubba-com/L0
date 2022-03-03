CREATE TABLE orders (
    order_uid varchar(128) PRIMARY KEY NOT NULL UNIQUE,
    track_number       varchar(128)        NOT NULL,
    entry              varchar(128)        NOT NULL,
    locale             varchar(128)        NOT NULL,
    internal_signature varchar(128)        NOT NULL,
    customer_id        varchar(128)        NOT NULL,
    delivery_service   varchar(128)        NOT NULL,
    shardkey           varchar(128)        NOT NULL,
    sm_id              bigint,
    date_created       timestamp           NOT NULL,
    oof_shard          varchar(128)        NOT NULL
);

SELECT
       orders.order_uid,
       orders.track_number,
       orders.entry,
       deliveries.order_uid,
       deliveries.name,
       deliveries.phone,
       deliveries.zip,
       deliveries.city,
       deliveries.address,
       deliveries.region,
       deliveries.email,
       payments.transaction,
       payments.request_id,
       payments.currency,
       payments.provider,
       payments.amount,
       payments.payment_dt,
       payments.bank,
       payments.delivery_cost,
       payments.goods_total,
       payments.custom_fee,
       orders.locale,
       orders.internal_signature,
       orders.customer_id,
       orders.delivery_service,
       orders.shardkey,
       orders.sm_id,
       orders.date_created,
       orders.oof_shard
FROM orders WHERE order_uid = $1
JOIN deliveries ON deliveries.order_uid = orders.order_uid
JOIN payments ON payments.order_uid = orders.order_uid;