# Configuration Documentation

This document describes how to configure the SAQ Inventory System Backend.

## Configuration Loading
The backend uses environment variables for configuration. When running the server locally, it reads environment variables from a `.env` file in the root directory using the `joho/godotenv` package.

---

## Environment Variables

| Variable | Description | Default | Required |
|----------|-------------|---------|----------|
| `DB_PATH` | File path to the SQLite database. | None | Yes |
| `STORAGE_PATH` | Directory where uploaded files and images are stored. | `./storage` | No |
| `PORT` | TCP port the HTTP server listens on. | `8080` | No |

---

## Example `.env` File
To start quickly, copy the provided `.env.example` to `.env` and adjust the values:

```ini
DB_PATH=database/saq_inventory.db
STORAGE_PATH=./storage
PORT=8080
```

---

## Production Configurations
In a production containerized environment (e.g., Docker), these values should be injected as container environment variables rather than using a physical `.env` file:

```yaml
environment:
  - DB_PATH=/var/lib/saq/saq_inventory.db
  - STORAGE_PATH=/var/lib/saq/storage
  - PORT=8080
```
Ensure that the directory housing the SQLite database (`DB_PATH`) and the `STORAGE_PATH` are mapped to persistent volumes so data is not lost on container restart.
