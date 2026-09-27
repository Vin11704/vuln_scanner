# Setup and Run

## Prerequisites

- Go 1.27 or later

## Setup

```powershell
git clone https://github.com/Vin11704/vuln_scanner.git
cd vuln_scanner
go mod download
```

Place the codebase you want to scan next to the scanner directory e.g.:

```text
../<codebase_path>
```

The target codebase should contain a `.gitignore` file.

## Tests

Run all tests with:

```powershell
go test ./...
```

To include detailed test output, use:

```powershell
go test -v ./...
```

Run one test by name from the current package:

```powershell
go test -run '^TestScan$'
```

## Run

Run directly with Go:

```powershell
go run .
```

Or build and run the executable:

```powershell
go build -o vuln_scanner.exe .
.\vuln_scanner.exe
```

# WIP

A simple vulnerability scanner I made for codebases. Currently only checks for leaked secrets that are not gitignored. Will add more functionalities in the future.
