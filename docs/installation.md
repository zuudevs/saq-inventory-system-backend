# Installation Documentation

This document describes how to install prerequisites, configure the project, run database migrations, and start the application.

## Prerequisites
Ensure the following tools are installed on your machine:
* Go 1.25 or higher
* SQLite 3
* Goose (Go Database Migrations tool)
* PowerShell (for Windows script execution)

---

## 1. Automated Script Installation (Windows PowerShell)

A PowerShell script is provided to automate installation of Go, SQLite, Goose, and the project's dependencies:

1. Open PowerShell with Administrator privileges.
2. Navigate to the project root directory.
3. Run the setup script:
   ```powershell
   .\scripts\setup.ps1
   ```

This script will:
* Check for and install Go using `winget`.
* Check for and install Goose using `go install`.
* Check for and install SQLite CLI using `winget`.
* Run `go mod download` to fetch project dependencies.

---

## 2. Manual Installation

If you prefer to install packages manually:

### Step 1: Install Go
Download and install Go (>= 1.25) from the official website.

### Step 2: Install Goose
Install Goose using Go:
```bash
go install github.com/pressly/goose/v3/cmd/goose@v3.24.1
```
Make sure `$GOPATH/bin` (or `%USERPROFILE%\go\bin` on Windows) is in your system's PATH.

### Step 3: Install SQLite
Install SQLite CLI.
* **On macOS (via Homebrew)**: `brew install sqlite`
* **On Ubuntu/Debian**: `sudo apt install sqlite3`
* **On Windows (via winget)**: `winget install SQLite.SQLite`

---

## 3. Configuration Setup
Copy the example environment file to `.env`:

* **Bash (macOS/Linux)**:
  ```bash
  cp .env.example .env
  ```
* **PowerShell (Windows)**:
  ```powershell
  Copy-Item .env.example .env
  ```

Ensure you edit `.env` and set `DB_PATH` to your desired database path (e.g. `database/saq_inventory.db`).

---

## 4. Run Migrations & Start App (Automated)

The project includes an autorun script that automatically loads environment variables, creates the database directories, runs goose migrations, and boots the backend server.

Run the autorun script:
```powershell
.\scripts\autorun.ps1
```

---

## 5. Run Migrations & Start App (Manual)

If running manually, execute:

### Step 1: Create Database Directory
```bash
mkdir -p database
```

### Step 2: Run Database Migrations
Run goose migrations against the SQLite database:
```bash
goose -dir migrations sqlite database/saq_inventory.db up
```

### Step 3: Start the Backend Server
```bash
go run cmd/server/main.go
```
The server will start listening on the port configured in `.env` (default is 8080).
