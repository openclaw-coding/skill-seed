# Build Instructions

This document explains how to build grow-check with different language settings.

## Build English Version (Default)

To build the English version (default):

```bash
go build -o grow-check ./cmd/grow-check
```

Or explicitly without language tags:

```bash
go build -tags="" -o grow-check ./cmd/grow-check
```

## Build Chinese Version

To build the Chinese version:

```bash
go build -tags cn -o grow-check-cn ./cmd/grow-check
```

## Installation

After building, you can install the binary to your system:

```bash
# For English version
mv grow-check /usr/local/bin/

# For Chinese version
mv grow-check-cn /usr/local/bin/grow-check
```

Or use `go install`:

```bash
# English version
go install ./cmd/grow-check

# Chinese version
go install -tags cn ./cmd/grow-check
```

## Development

When developing, you can test both language versions:

```bash
# Test English
go run ./cmd/grow-check check

# Test Chinese
go run -tags cn ./cmd/grow-check check
```

## Language Files

- `internal/i18n/i18n.go` - English messages (default)
- `internal/i18n/i18n_cn.go` - Chinese messages (requires `cn` build tag)

## Adding New Languages

To add a new language:

1. Create a new file in `internal/i18n/` (e.g., `i18n_es.go` for Spanish)
2. Add the appropriate build tag at the top:
   ```go
   //go:build es
   // +build es
   ```
3. Copy the message map and translate all values
4. Build with: `go build -tags es ./cmd/grow-check`
