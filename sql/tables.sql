CREATE DATABASE packagetracer ENCODING 'UTF8';

CREATE TABLE IF NOT EXISTS customs (
    refCode VARCHAR(16) NOT NULL,
    status VARCHAR(128) NOT NULL,
    detail VARCHAR(256) NOT NULL,
    date TIMESTAMP NOT NULL,
    PRIMARY KEY (refCode, date),
    UNIQUE (refCode, date)
);