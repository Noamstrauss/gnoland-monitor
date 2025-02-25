
-include .env

# Default configs
SUB_ACCOUNTS ?= 5
TRANSACTIONS ?= 100
MNEMONIC ?= source bonus chronic canvas draft south burst lottery vacant surface solve popular case indicate oppose farm nothing bullet exhibit title speed wink action roast
OUTPUT_FILE ?= result.json
NODE_PORT ?= 26657
INDEXER_PORT ?= 8546
METRICS_PORT ?= 8080
GRAFANA_PORT ?= 3000
PROMETHEUS_PORT ?= 9090
PROCESSING_INTERVAL ?= 5s
DB_PATH ?= /data/indexer-db
PROMETHEUS_VERSION ?= v3.2.0
GRAFANA_VERSION ?= 11.5.2

export

.PHONY: build up down restart logs status clean help

env-setup:
	@if [ ! -f .env ]; then \
		echo "Creating .env file from example..."; \
		cp .env.example .env; \
		echo ".env file created. You can now edit it with your custom settings."; \
	else \
		echo ".env file already exists."; \
	fi

build:
	@echo "Building Docker images..."
	docker-compose build

up: env-setup
	@echo "Starting services with configuration:"
	@echo "- SUB_ACCOUNTS: $(SUB_ACCOUNTS)"
	@echo "- TRANSACTIONS: $(TRANSACTIONS)"
	@echo "- OUTPUT_FILE: $(OUTPUT_FILE)"
	@echo "- MNEMONIC: $(MNEMONIC)"
	#docker-compose up -d
	envsubst < prometheus/prometheus.yml.template > prometheus/prometheus.yml
	envsubst < grafana/provisioning/datasources/datasources.yml.template > grafana/provisioning/datasources/datasources.yml
	docker-compose up --build -d
	@echo "Services started:"
	@echo "- Gnoland: http://localhost:$(NODE_PORT)"
	@echo "- Indexer: http://localhost:$(INDEXER_PORT)/graphql"
	@echo "- Metrics: http://localhost:$(METRICS_PORT)/metrics"
	@echo "- Grafana: http://localhost:$(GRAFANA_PORT)/d/gno-node-monitoring"
	@echo "- Prometheus: http://localhost:$(PROMETHEUS_PORT)"


down:
	@echo "Stopping services..."
	docker-compose down

restart: down up

logs:
	@if [ "$(service)" ]; then \
		docker-compose logs -f $(service); \
	else \
		docker-compose logs -f; \
	fi

status:
	docker-compose ps

clean:
	docker-compose down -v
	rm -rf ./results/*

help:
	@echo "Available commands:"
	@echo "  make build                    - Build all Docker images"
	@echo "  make up                       - Start all services with default config"
	@echo "  make up SUB_ACCOUNTS=10 ...   - Start with custom config"
	@echo "  make down                     - Stop all services"
	@echo "  make restart                  - Restart all services"
	@echo "  make logs [service=name]      - View logs"
	@echo "  make status                   - Check service status"
	@echo "  make clean                    - Clean up everything"
	@echo ""
	@echo "Current configuration:"
	@echo "  SUB_ACCOUNTS = $(SUB_ACCOUNTS)"
	@echo "  TRANSACTIONS = $(TRANSACTIONS)"
	@echo "  OUTPUT_FILE  = $(OUTPUT_FILE)"
	@echo "  MNEMONIC     = $(MNEMONIC)"
	@echo "  GNOLAND 	  = http://localhost:$(NODE_PORT)"
	@echo "  INDEXER      = http://localhost:$(INDEXER_PORT)/graphql"
	@echo "  METRICS:     = http://localhost:$(METRICS_PORT)/metrics"
	@echo "  GRAFANA:     = http://localhost:$(GRAFANA_PORT)"
	@echo "  PROMETHEUS:  = http://localhost:$(PROMETHEUS_PORT)"