# Changelog

All notable changes to this project will be documented in this file.

---

## [Unreleased]

### Added
* XLSX Import endpoint `POST /imports/xlsx` for importing system data from Excel workbooks.
* Strict XLSX validation engine checking required sheet names (`Brands`, `Categories`, `Locations`, `Items`, `Images`), column headers (Row 1), and data cell types across all resources.
* Atomic database transaction support for import execution, rolling back upon any validation or constraint error.
* Export endpoints `/exports/csv` and `/exports/xlsx` for downloading all system resources (`Brands`, `Categories`, `Locations`, `Items`, `Images`).
* ZIP archive bundling (`archive/zip`) for CSV exports, generating individual CSV files per resource (`brands.csv`, `categories.csv`, etc.).
* Multi-sheet Excel workbook export (`ExportMultiSheetXLSX`) generating distinct worksheets for each resource.
* Dynamic schema alter support using `ALTER TABLE` (`ADD COLUMN` and `DROP COLUMN`) when updating a category's metadata structure, preventing data loss in unchanged columns.
* HTTP endpoint `PUT /categories/{categoryId}/metadata-structure` to update metadata structures dynamically.
* HTTP endpoint `DELETE /categories/{categoryId}/metadata-structure` to safely delete metadata structures and drop their physical SQLite tables.
* Restructured integration test suite into modular scripts under `tests/api/` and updated `tests/api_test.ps1`.

### Changed
* Refactored export endpoints from `/exports/items/*` to `/exports/csv` and `/exports/xlsx`.
* Expanded export scope from item-only resources to all system resources.

### Fixed
* Category validation check bug (`category == nil`) in `MetadataStructureService.Update`.
* Fixed database update bug where `structure.ID` was not mapped prior to updating metadata structures.

---

## [v1.0.0-sqlite-alpha]

### Added
* Comprehensive repository documentation including API, architecture, configuration, database schema, deployment, development, and installation guides.
* Project configuration and compliance files: Code of Conduct, Contribution guidelines, and Security policies.
* Multipart image upload functionality under `/images/upload` using local file system storage.
* Automated integration test script (`tests/api_test.ps1`) using native curl execution.
* Setup script (`scripts/setup.ps1`) for downloading Go, Goose, and SQLite CLI tools automatically.
* Support for SQLite database migration workflows via CI/CD pipelines.

### Changed
* Migrated primary database from MySQL to SQLite (modernc.org driver) to simplify setup.
* Enhanced metadata inclusion directly inside item listing and item retrieval API responses.
* Optimized docker-compose configurations and automated launch scripts.

---

## [v1.0.0-sqlite-pre-alpha]

### Added
* Support for dynamic metadata structures per category, enabling customizable attributes for inventory items.
* Automatic SQLite triggers to handle updating timestamps (`updated_at` triggers).
* Foreign key constraints mapping brands, categories, locations, and items.
* Integration of the `zuu-powershell-dotenv` tool for loading environmental files automatically in PowerShell.

### Changed
* Refactored repository layer query execution, ensuring complete entity payloads are returned following insertion or update operations.
* Formatted source code to strictly align with Go code formatting guidelines.

---

## [mysql-v1.0.0-pre-alpha]

### Added
* Initial draft of the backend implementation using a MySQL database setup.
* Basic handlers and routing bindings for Brand, Category, Item, and Location CRUD APIs.
* Containerized multi-stage Docker build pipeline configuration.
