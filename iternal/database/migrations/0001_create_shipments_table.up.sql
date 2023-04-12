CREATE TABLE shipments
(
    id int primary key generated always as identity,
    barcode varchar(25) not null unique,
    sender varchar(30) not null,
    receiver varchar(30) not null,
    is_delivered boolean default false,
    origin varchar(30) not null,
    destination varchar(30) not null,
    created_at timestamp default now() not null
);

CREATE INDEX idx_barcode ON shipments(barcode);
CREATE INDEX idx_id on shipments(id);
