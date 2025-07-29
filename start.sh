#!/bin/sh
set -e
# single-node, in-memory, no TLS for local dev
cockroach start-single-node \
  --insecure \
  --listen-addr=localhost:26257 \
  --http-addr=localhost:8080 \
  --background

# wait for SQL port
until cockroach sql --insecure -e "SELECT 1" >/dev/null 2>&1; do sleep 1; done
# create DB/user if you wish
cockroach sql --insecure -e "CREATE DATABASE IF NOT EXISTS myapp"

# launch nginx (your Go backend is already compiled in)
nginx -g "daemon off;"