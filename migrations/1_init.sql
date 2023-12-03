-- +migrate Up
CREATE SCHEMA IF NOT EXISTS testdata;
-- Drop
DROP TABLE IF EXISTS testdata.stocks;
-- stock
CREATE TABLE testdata.stocks
(
    stock_id     SERIAL,
    stock_name   Char(20) not null, --limit at 20 char
    stock_price  int, -- small calculate at test
    last_update  timestamp,
    PRIMARY KEY (stock_id)
);

INSERT INTO testdata.stocks(stock_name, stock_price, last_update)
VALUES ('stock_1', 111, current_timestamp);
INSERT INTO testdata.stocks(stock_name, stock_price, last_update)
VALUES ('stock_2', 222, current_timestamp);
INSERT INTO testdata.stocks(stock_name, stock_price, last_update)
VALUES ('stock_3', 333, current_timestamp);
INSERT INTO testdata.stocks(stock_name, stock_price, last_update)
VALUES ('stock_4', 444, current_timestamp);
INSERT INTO testdata.stocks(stock_name, stock_price, last_update)
VALUES ('stock_5', 5555, current_timestamp);
INSERT INTO testdata.stocks(stock_name, stock_price, last_update)
VALUES ('stock_6', 66666, current_timestamp);