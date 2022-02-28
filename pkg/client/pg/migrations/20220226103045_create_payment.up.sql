CREATE TABLE payments (
    transaction   varchar PRIMARY KEY NOT NULL UNIQUE,
    request_id    varchar,
    currency      varchar(128)        NOT NULL,
    provider      varchar(128)        NOT NULL,
    amount        bigint              NOT NULL,
    payment_dt    bigint              NOT NULL,
    bank          varchar(128)        NOT NULL,
    delivery_cost bigint              NOT NULL,
    goods_total   bigint              NOT NULL,
    custom_fee    bigint              NOT NULL,

    order_uid varchar NOT NULL UNIQUE,
    FOREIGN KEY (order_uid) REFERENCES orders (order_uid) ON DELETE CASCADE
)