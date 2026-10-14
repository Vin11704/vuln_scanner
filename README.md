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
### Testing steps 
1. Formats the Go test file in place.
```powershell
gofmt -w .\your_new_test_file_test.go
```

2. Runs only tests related to the specified features.
```powershell
go test -run '^Test(YourFeature1|YourFeature2)' .
```

3. Runs all tests in every Go package in the repository.
```powershell
go test d:.
```

4. Performs static analysis to find suspicious or likely buggy Go code. It does not run tests.
```powershell
go vet d:.
```     

<!-- 5. Runs all tests with Go’s race detector enabled to find unsafe concurrent memory access. It is slower but useful because the scanner uses workers and goroutines.
```powershell
go test -race .
``` -->


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

# Guide
## To create new detection feature:
- implement `Rule` from [rule.go](./rule.go). Ensure all functions in interface is used.
- use AppliesTo function to check if the rule applies to the file being scanned. For secrets rule, it returns True every time since secrets can be anywhere in the codebase. 
- [patterns.go](./patterns.go) for the regex stuff.
- `Check` function 

## TBC
- check gitignore functions in main.go
- check [gitignore.go](./gitignore.go)


# WIP

A simple vulnerability scanner I made for codebases. Currently only checks for leaked secrets that are not gitignored. Will add more functionalities in the future.

