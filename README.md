![Gnoland logo](https://gnolang.github.io/blog/2024-05-21_the-gnome/src/banner.png)

# Gnoland Monitor


![Lint](https://github.com/noamstrauss/gnoland-monitor/actions/workflows/lint.yaml/badge.svg)

A monitoring stack for tracking metrics and performance of Gnoland using Grafana, Prometheus, and a custom metrics processor.


## Table of Contents

- [Pre-requisites](#Pre-requisites)
- [Features](#features)
- [Quick Start](#quick-start)
- [Configuration](#configuration)
- - [Environment Variables](#environment-variables)
- - [Using .env File](#using-env-file)
- [Usage](#usage)
- - [Common Commands](#common-commands)
- - [Accessing Services](#accessing-services)
- [Development](#development)
- [Troubleshooting](#troubleshooting)
- [Contributing](#contributing)
- [License](#license)


## Pre-requisites

Before using this GnoLand Monitor, ensure that you have the following:

- [Go](https://go.dev/doc/install) (v1.19 or later)
- [Docker](https://docs.docker.com/engine/install/) (v20.10 or later)
- [Docker Compose](https://docs.docker.com/compose/install) (v2.0 or later)
- [Make](https://formulae.brew.sh/formula/make)
- [Git](https://git-scm.com/downloads)


## Features

* **Real-time Metrics:** Collect and visualize node performance metrics in real-time
* **Grafana:** Pre-configured Grafana dashboard for comprehensive monitoring
* **Transaction Indexing:** GraphQL-based transaction indexer for easy querying and analysis
* **Supernova Integration:** Built-in support for Supernova node for network interaction
* **Dockerized Setup:** Easy deployment and management using Docker Compose

# Quick Start

**Clone the repository:**

```bash
git clone https://github.com/noamstrauss/gnoland-monitor.git
cd gnoland-monitor
```

**Configure environment variables (optional):**
```bash
cp .env.example .env
# Edit .env file with your preferred settings
```

**Start the services:**

```bash
make up
```

**Access Grafana Dashboard:**

```
URL: http://localhost:3000/d/gno-node-monitoring
```

**Stop the services:**

```bash
make down
````

# Components

Gnoland Monitor consists of the following components:

* **Gnoland Node:** Core Gnoland blockchain node
* **Transaction Indexer:** GraphQL service for indexing and querying transactions
* **Metrics Processor:** Custom Go service that processes and exposes Prometheus metrics
* **Supernova:** Tool for generating test transactions and network load
* **Prometheus:** Time-series database for storing metrics
* **Grafana:** Visualization platform with pre-configured dashboards


# Configuration

## Environment Variables

You can customize the behavior of Gnoland Monitor using the following environment variables:

| Variable            | Description                          | Default Value                          |
|---------------------|----------------------------------|----------------------------------|
| SUB_ACCOUNTS       | Number of sub-accounts to monitor | 5                                |
| TRANSACTIONS       | Number of transactions to generate | 100                              |
| MNEMONIC           | Mnemonic for account generation  | source bonus chronic...         |
| OUTPUT_FILE        | Output file for transaction results | result.json                     |
| NODE_PORT          | Gnoland node port                | 26657                            |
| INDEXER_PORT       | Transaction indexer port        | 8546                             |
| METRICS_PORT       | Metrics processor port         | 8080                             |
| GRAFANA_PORT       | Grafana dashboard port         | 3000                             |
| PROMETHEUS_PORT    | Prometheus port               | 9090                             |
| PROCESSING_INTERVAL | Metrics processing interval   | 5s                               |
| DB_PATH            | Path to indexer database      | /data/indexer-db                 |
| PROMETHEUS_VERSION | Prometheus version           | v3.2.0                           |
| GRAFANA_VERSION    | Grafana version              | 11.5.2                           |


## Using .env File

For easier configuration, you can create a `.env` file in the project root with your desired settings:

```bash
#Node Configuration
SUB_ACCOUNTS=10
TRANSACTIONS=200
MNEMONIC=your custom mnemonic here
OUTPUT_FILE=custom-result.json

# Service Ports
NODE_PORT=26657
INDEXER_PORT=8546
METRICS_PORT=8080
GRAFANA_PORT=3000
PROMETHEUS_PORT=9090

# Other Settings
PROCESSING_INTERVAL=10s
DB_PATH=/data/indexer-db
PROMETHEUS_VERSION=v3.2.0
GRAFANA_VERSION=11.5.2
```

# Usage

## Common Commands

```bash
  make build                    - Build all Docker images
  make up                       - Start all services with default config
  make up SUB_ACCOUNTS=10 ...   - Start with custom config
  make down                     - Stop all services
  make restart                  - Restart all services
  make logs [service=name]      - View logs
  make status                   - Check service status
  make clean                    - Clean up everything
```

## Accessing Services

After starting the services, you can access them at the following URLs:

* **Gnoland Node:** http://localhost:26657
* **Transaction Indexer:** http://localhost:8546/graphql
* **Metrics Endpoint:** http://localhost:8080/metrics
* **Grafana Dashboard:** http://localhost:3000/d/gno-node-monitoring
* **Prometheus UI:** http://localhost:9090

# Development

To build and run the metrics-processor service locally:

```bash
cd metrics-processor
go build -o metrics-processor
./metrics-processor
```

To modify the Grafana dashboards:

* Access the Grafana UI at http://localhost:3000
* Make your changes to the dashboard
* Export the dashboard to JSON
* Save the JSON to `grafana/provisioning/dashboards/gnoland-dashboard.json`
* Rebuild and restart the services: make restart

# Troubleshooting

Common issues and their solutions:

1. **Services failing to start**: Check logs with make logs to identify the specific issue
2. **Cannot connect to Grafana**: Ensure ports are not in use by other applications
3. **No metrics showing up**: Verify that the metrics-processor is healthy with make status
4. **Database errors**: Try make clean to reset all data volumes

# Contributing

Contributions are welcome! Please follow these guidelines:

1. Fork the repository.
2. Create a new branch for your feature or bug fix.
3. Follow the project's coding style.
4. Submit a pull request with a clear description of your changes.

# License

This project is licensed under the MIT License. See the LICENSE file for details. 

