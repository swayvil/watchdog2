#!/bin/bash
psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" -d "$POSTGRES_DB" <<-EOSQL
    CREATE USER watchdog WITH PASSWORD 'watchdog';
    CREATE DATABASE watchdog2 WITH OWNER watchdog;
EOSQL
psql -v ON_ERROR_STOP=1 --username "watchdog" -d "watchdog2" <<-EOSQL
    CREATE TABLE snapshot(
        id SERIAL PRIMARY KEY,
        camera CHAR(10)         NOT NULL,
        timestamp TIMESTAMP     NOT NULL,
        internaltimestamp TIMESTAMP     NOT NULL
    );
    CREATE TABLE camera (
        camera CHAR(10) PRIMARY KEY         NOT NULL
    );
    INSERT INTO camera VALUES ('Cour');
    INSERT INTO camera VALUES ('Garage');
    INSERT INTO camera VALUES ('Entree');
EOSQL