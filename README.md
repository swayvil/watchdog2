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
    timestamp TIMESTAMP     NOT NULL
);

CREATE TABLE camera (
    camera CHAR(10) PRIMARY KEY         NOT NULL
);

INSERT INTO camera (camera) VALUES ('Cour');
INSERT INTO camera (camera) VALUES ('Garage');
INSERT INTO camera (camera) VALUES ('Entree');

mkdir -p /Users/idaho/workspaces/watchdog2-store/cour/small
mkdir -p /Users/idaho/workspaces/watchdog2-store/cour/big
mkdir -p /Users/idaho/workspaces/watchdog2-store/garage/small
mkdir -p /Users/idaho/workspaces/watchdog2-store/garage/big
mkdir -p /Users/idaho/workspaces/watchdog2-store/entree/small
mkdir -p /Users/idaho/workspaces/watchdog2-store/entree/big

DROP TABLE snapshot;
DROP TABLE camera;
rm -Rf /Users/idaho/workspaces/watchdog2-store

https://godoc.org/github.com/lib/pq#Open
https://www.calhoun.io/connecting-to-a-postgresql-database-with-gos-database-sql-package/
https://github.com/emersion/go-imap

go get github.com/lib/pq
go get github.com/gorilla/mux
go get github.com/rs/cors
go build


http://localhost:9999/snapshots-all-cams/2019-11-01T05:41:00


curl -i -H "Accept: application/json" -H "Content-Type: application/json" -X GET http://localhost:9999/snapshots-all-cams/2019-11-01T05:41:00
