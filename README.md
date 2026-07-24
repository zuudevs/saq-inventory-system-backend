# SAQ Inventory System Backend

> A backend system for managing inventory items, brands, categories, locations, and dynamic metadata.

## Features

- CRUD operations for Brands, Categories, Locations, and Items
- Image upload and management with static file serving
- Dynamic metadata structure for categories
- Full-system data export to CSV (ZIP archive) and XLSX (multi-sheet workbook)
- SQLite database integration with foreign key support
- Automated API testing
- Docker support

## Tech Stack

- Go 1.25
- SQLite
- Docker
- Chi Router
- sqlx

## Requirements

- Go >= 1.25
- Docker (optional)
- SQLite

## Installation

```bash
git clone https://github.com/zuudevs/saq-inventory-system-backend.git
cd saq-inventory-system-backend

go mod download
```

## Configuration

```bash
cp .env.example .env
```

Edit `.env` according to your environment.

## Running

```bash
go run cmd/main.go
```

or

```bash
docker compose up
```

## Testing

```bash
go test ./...
```

## Project Structure

```
cmd/          # Entry point of the application
internal/     # Core business logic and internal packages
migrations/   # Database migration files
tests/        # API and unit tests
database/     # Database storage
scripts/      # Utility scripts
tools/        # Development tools
```

## API Documentation

See [`docs/`](docs/).

## Contributing

1. Fork repository
2. Create a feature branch
3. Commit changes
4. Open a Pull Request

## License

MIT