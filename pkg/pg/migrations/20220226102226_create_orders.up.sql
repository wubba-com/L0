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