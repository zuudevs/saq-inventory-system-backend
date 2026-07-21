# Deployment Documentation

This document describes how to deploy the SAQ Inventory System Backend to production environments.

## Build and Run Native Go Binary
To compile the project into a single, high-performance static binary:

```bash
# Compile
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o server ./cmd/server

# Set production env parameters
export DB_PATH=/var/lib/saq/saq_inventory.db
export STORAGE_PATH=/var/lib/saq/storage
export PORT=8080

# Run
./server
```

---

## Containerized Deployment (Docker)

A multi-stage `Dockerfile` is provided for containerizing the application. It builds the server binary in a Go environment, then packs it inside a lightweight Alpine image.

### Building the Docker Image
```bash
docker build -t saq-inventory-backend:latest .
```

### Running the Container with Volume Mounts
Since the application uses an SQLite database and stores uploaded images locally on the file system, you must mount persistent volumes to avoid data loss.

```bash
docker run -d \
  --name saq-backend \
  -p 8080:8080 \
  -v /var/lib/saq/database:/app/database \
  -v /var/lib/saq/storage:/app/storage \
  -e DB_PATH=database/saq_inventory.db \
  -e STORAGE_PATH=./storage \
  -e PORT=8080 \
  saq-inventory-backend:latest
```

---

## Docker Compose Setup

*Note: The included `docker-compose.yml` file contains a MySQL service and backend database variables (`DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASS`, `DB_NAME`). However, the current Go codebase is compiled to use SQLite via `DB_PATH`.*

To run SQLite under Docker Compose, update your configuration or inject the `DB_PATH` variable:

```yaml
version: "3.9"

services:
  backend:
    build: .
    restart: unless-stopped
    environment:
      - DB_PATH=/app/database/saq_inventory.db
      - STORAGE_PATH=/app/storage
      - PORT=8080
    ports:
      - "8080:8080"
    volumes:
      - saq_db:/app/database
      - saq_storage:/app/storage

volumes:
  saq_db:
  saq_storage:
```
