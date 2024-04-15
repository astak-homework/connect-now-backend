#!/bin/bash
set -e
postgres_replicator_password=`cat /run/secrets/postgres_replicator_password`
psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
    CREATE USER replicator WITH REPLICATION ENCRYPTED PASSWORD '${postgres_replicator_password}';
    SELECT pg_create_physical_replication_slot('replication_slot');
EOSQL