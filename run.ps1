$nats = Start-Process -PassThru -FilePath "output\nats\nats-server.exe" -ArgumentList "--addr localhost --port 44222 --server_name `"My Server`"" -WorkingDirectory "output\nats"
$backend_service = Start-Process -PassThru -FilePath "output\backend-service.exe" -ArgumentList "-nats_url nats://localhost:44222" -WorkingDirectory "output"
$frontend_service = Start-Process -PassThru -FilePath "output\frontend-service.exe" -ArgumentList "-nats_url nats://localhost:44222 -port 8090" -WorkingDirectory "output"
$caddy = Start-Process -PassThru -FilePath "output\caddy\caddy.exe" -ArgumentList "run --config Caddyfile.local" -WorkingDirectory "."

echo "Press any key to stop ..."
[Console]::ReadKey()

Stop-Process -InputObject $caddy
Stop-Process -InputObject $frontend_service
Stop-Process -InputObject $backend_service
Stop-Process -InputObject $nats