CREATE TABLE deliveries (
    order_uid varchar      NOT NULL,
    name      varchar(128) NOT NULL,
    phone     varchar(16)  NOT NULL,
    zip       varchar(128) NOT NULL,
    city      varchar(128) NOT NULL,
    address   varchar(256) NOT NULL,
    region    varchar(256) NOT NULL,
    email     varchar(128) NOT NULL,

    UNIQUE (order_uid),
    CONSTRAINT fk_order FOREIGN KEY (order_uid) REFERENCES orders (order_uid) ON DELETE CASCADE
);