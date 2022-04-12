CADDY_WINDOWS_URL := https://github.com/caddyserver/caddy/releases/download/v2.4.6/caddy_2.4.6_windows_amd64.zip
CADDY_LINUX_URL := https://github.com/caddyserver/caddy/releases/download/v2.4.6/caddy_2.4.6_linux_amd64.tar.gz
NATS_WINDOWS_URL := https://github.com/nats-io/nats-server/releases/download/v2.7.4/nats-server-v2.7.4-windows-amd64.zip
NATS_LINUX_URL := https://github.com/nats-io/nats-server/releases/download/v2.7.4/nats-server-v2.7.4-linux-amd64.tar.gz

output_folder := output

ifeq ($(OS), Windows_NT)
	BINARY_SUFFIX := .exe
	CADDY_ARCHIVE_NAME := $(notdir $(CADDY_WINDOWS_URL))
	DOWNLOAD_CADDY_CMD := powershell -Command "Invoke-WebRequest $(CADDY_WINDOWS_URL) -OutFile $(CADDY_ARCHIVE_NAME)"
	EXTRACT_CADDY_CMD := powershell -Command "Expand-Archive -Force $(CADDY_ARCHIVE_NAME) ."
	NATS_ARCHIVE_NAME := $(notdir $(NATS_WINDOWS_URL))
	DOWNLOAD_NATS_CMD := powershell -Command "Invoke-WebRequest $(NATS_WINDOWS_URL) -OutFile $(NATS_ARCHIVE_NAME)"
	EXTRACT_NATS_CMD := powershell -Command "Expand-Archive -Force $(NATS_ARCHIVE_NAME) ."
else
	BINARY_SUFFIX := 
	CADDY_ARCHIVE_NAME := $(notdir $(CADDY_LINUX_URL))
	DOWNLOAD_CADDY_CMD := powershell -Command "Invoke-WebRequest $(CADDY_LINUX_URL) -OutFile $(CADDY_ARCHIVE_NAME)"
	EXTRACT_CADDY_CMD := powershell -Command "Expand-Archive -Force $(CADDY_ARCHIVE_NAME) ."
	EXTRACT_NATS_CMD := tar -xf $(NATS_ARCHIVE_NAME)
endif

.PHONY: application

build: backend_service frontend_service application

backend_service:
	cd backend-service && \
	go build -o $(addsuffix $(BINARY_SUFFIX), ../$(output_folder)/backend-service)
	
frontend_service:
	cd frontend-service && \
	go build -o $(addsuffix $(BINARY_SUFFIX), ../$(output_folder)/frontend-service)

create_output_folder:
	mkdir $(output_folder) || \
	true

create_caddy_folder:
	cd $(output_folder) && \
	mkdir caddy || \
	true
	
create_nats_folder:
	cd $(output_folder) && \
	mkdir nats || \
	true

download: caddy nats

caddy: create_output_folder create_caddy_folder
	cd $(output_folder)/caddy && \
	$(DOWNLOAD_CADDY_CMD) && \
	$(EXTRACT_CADDY_CMD) && \
	rm $(CADDY_ARCHIVE_NAME)
	
nats: create_output_folder create_nats_folder
	cd $(output_folder)/nats && \
	$(DOWNLOAD_NATS_CMD) && \
	$(EXTRACT_NATS_CMD) && \
	mv $(basename $(NATS_ARCHIVE_NAME))/* . && \
	rm -r $(basename $(NATS_ARCHIVE_NAME)) && \
	rm $(NATS_ARCHIVE_NAME)
	
application:
	cd application && \
	npm run build