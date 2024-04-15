#!/bin/bash
set -e
pg_basebackup -h pgmaster -D /var/lib/postgresql/data --wal-method=stream -R --slot=replication_slot