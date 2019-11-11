# watchdog2

docker pull postgres:latest
docker create --name postgres \
-e POSTGRES_PASSWORD=watchdog2 \
-p 5432:5432 \
--volume postgres:/var/lib/postgresql/data \
postgres

docker exec -ti postgres /bin/bash

psql -U postgres

CREATE USER watchdog WITH PASSWORD 'watchdog';

CREATE DATABASE watchdog2 WITH OWNER watchdog;

psql -U watchdog watchdog2

CREATE TABLE snapshot(
    id SERIAL PRIMARY KEY,
    camera CHAR(10)         NOT NULL,
    timestamp TIMESTAMP     NOT NULL,
    photosmall INT          NOT NULL,
    photo INT               NOT NULL
);

CREATE TABLE photosmall(
    id SERIAL PRIMARY KEY,
    photo BYTEA             NOT NULL
);

CREATE TABLE photo(
    id SERIAL PRIMARY KEY,
    photo BYTEA             NOT NULL
);

CREATE TABLE camera (
    camera CHAR(10) PRIMARY KEY         NOT NULL
);

INSERT INTO camera (camera) VALUES ('cours');
INSERT INTO camera (camera) VALUES ('garage');
INSERT INTO camera (camera) VALUES ('entree');


https://godoc.org/github.com/lib/pq#Open
https://www.calhoun.io/connecting-to-a-postgresql-database-with-gos-database-sql-package/
https://github.com/emersion/go-imap

go get github.com/lib/pq
go build