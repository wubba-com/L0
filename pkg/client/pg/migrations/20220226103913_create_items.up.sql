CREATE TABLE items (
    chrt_id      bigint PRIMARY KEY NOT NULL,
    track_number varchar(256)       NOT NULL,
    price        bigint             NOT NULL,
    rid          varchar,
    name         varchar(128)       NOT NULL,
    sale         bigint             NOT NULL,
    size         varchar            NOT NULL,
    total_price  bigint             not null,
    nm_id        bigint,
    brand        varchar(256)       NOT NULL,
    status       int                NOT NULL,

    order_uid    varchar            NOT NULL,
    FOREIGN KEY (order_uid) REFERENCES orders (order_uid)
);