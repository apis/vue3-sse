Start-Process -FilePath "output\nats\nats-server.exe" -ArgumentList "--addr localhost --port 44222 --server_name `"My Server`"" -WorkingDirectory "output\nats"
Start-Process -FilePath "output\backend-service.exe" -ArgumentList "-nats_url nats://localhost:44222" -WorkingDirectory "output"
Start-Process -FilePath "output\frontend-service.exe" -ArgumentList "-nats_url nats://localhost:44222 -port 8090" -WorkingDirectory "output"
Start-Process -FilePath "output\caddy\caddy.exe" -ArgumentList "run --config Caddyfile.local" -WorkingDirectory "."