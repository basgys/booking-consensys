ifndef CONFIG_URI
	# File
	export CONFIG_URI=file://${PWD}/config/development.toml
	# Consul
	# port : 8500
	# export CONFIG_URI=consul://localhost:8500/config/booking-api.toml?dc=dc1
endif
ifndef IP
	export IP=0.0.0.0
endif
ifndef HTTP_PORT
	export HTTP_PORT=8484
endif
ifndef HOSTNAME
	export HOSTNAME="booking-api"
endif

target:
		@make build

dev:
		@echo ""
		@echo "🚀 Boot from source code ${HOSTNAME}..."
		@echo "\tOpen http://${IP}:${HTTP_PORT}"
		@echo ""
		@go run main.go

build:
		@go build ./...
		@go build -o booking-api

test:
		@go test -v -race ./...
