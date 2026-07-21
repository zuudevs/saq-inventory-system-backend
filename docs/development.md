# Development Documentation

This document serves as a guide for developers working on the SAQ Inventory System Backend.

## Project Structure

* `cmd/server/`: Entry point of the application, boots the SQLite database connection, sets up dependency injection, and runs the HTTP server.
* `internal/config/`: Database connection initialization and SQLite configuration (such as enabling foreign keys).
* `internal/dto/`: Data Transfer Objects for request validation and response mapping.
* `internal/handlers/`: HTTP request handlers that parse path parameters, request bodies, and invoke corresponding business services.
* `internal/models/`: Definitions of core domain structures (Brand, Category, Item, Location, Image, MetadataStructure).
* `internal/repositories/`: Implementation of SQL database operations utilizing `sqlx`.
* `internal/routes/`: Chi router mappings connecting handlers to API endpoints.
* `internal/schema/`: The logic for validating SQL identifiers, building dynamic DDL statements, and executing metadata table creations.
* `internal/services/`: Business services layer implementing validation rules, database transaction boundaries, and orchestration.
* `internal/utils/`: Utilities for handling JSON responses, uploads, and file system tasks.
* `migrations/`: Database schema version migrations managed using Goose.
* `tests/`: Complete API integration test scripts.

---

## Guidelines for Adding New Features

### 1. Model Definition
If adding a new static resource, define the structure in `internal/models`. Maintain the DB schema and json tags correctly.

### 2. Database Migration
Create a new SQL migration file in the `migrations` folder. Name it sequentially following the pattern `00000X_description.sql`. Use Goose SQL directives:
```sql
-- +goose Up
-- SQL statements here

-- +goose Down
-- SQL statements here
```

### 3. Repository
Create a repository file in `internal/repositories` to handle database operations (Create, Find, Update, Delete) using `sqlx`.

### 4. Service
Define the business rules, validation logic, and transaction bounds in `internal/services`.

### 5. DTO
Create DTO structs in `internal/dto` for request input and response serialization. Maintain clean mapping functions (`ToModel`, `ToResponse`).

### 6. Handlers & Routes
Create the handler file in `internal/handlers`, write response serialization logic via `utils.JSON`, and register endpoints in `internal/routes/routes.go`.

---

## Database Dynamic Schema Development
The dynamic schema service is a core component. When introducing new dynamic metadata field types:
* Update `internal/models/metadata_structure.go` by adding the new type to `MetadataFieldType`.
* Update `internal/schema/type_mapper.go` to support mapping the new field type to an appropriate SQLite column type.
* Ensure all identifier validations are strictly done using `schema.ValidateIdentifier` to prevent SQL injection vulnerabilities.

---

## Running Integration Tests

Integration and end-to-end tests are written in PowerShell, executing raw `curl` requests to ensure standard client compatibility across systems.

### Run Tests
```powershell
.\tests\api_test.ps1
```

### Run Tests with Image Upload
```powershell
.\tests\api_test.ps1 -TestUpload
```

### Run Tests without Deleting Generated Data
```powershell
.\tests\api_test.ps1 -SkipCleanup
```
