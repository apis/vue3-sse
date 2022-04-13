#!/bin/sh
#

(cd "output/nats" && ./nats-server --addr localhost --port 44222 --server_name "My Server")&
NATS_PID=$!

(cd "output" && ./backend-service -nats_url nats://localhost:44222)&
BACKEND_SERVICE_PID=$!

(cd "output" && ./frontend-service -nats_url nats://localhost:44222 -port 8090)&
FRONTEND_SERVICE_PID=$!

(output/caddy/caddy run --config Caddyfile.local)&
CADDY_PID=$!

echo "Press any key to stop..."
read REPLY

kill $NATS_PID
kill $BACKEND_SERVICE_PID
kill $FRONTEND_SERVICE_PID
kill $CADDY_PID