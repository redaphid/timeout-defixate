# timeout-defixate

Time Out app break enforcer for ADHD hyperfocus management.

## Build

```bash
go build -o timeout-defixate .
```

## Run

```bash
./timeout-defixate
# Or with custom limits:
./timeout-defixate --lock-limit=3 --shutdown-limit=5
```

## Requirements

- macOS (uses unified logging and private screen lock API)
- Time Out app with "Log to Console" enabled (Advanced → Diagnostic options)
