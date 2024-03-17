# GO Conways Game of Life & Websockets (Sandbox)

## Development

### Watch

Install Air

```bash
go install github.com/cosmtrek/air@latest
```

Watch

```bash
air
```

### Test

```bash
go test -benchmem -bench BenchmarkNewLifeRandom ./life/
go test -benchmem -bench BenchmarkDrawImageDataUrl ./life/
go test -benchmem -bench BenchmarkLife ./life/
```

## Build release

Optimize size:

- -s Omit the symbol table and debug information.
- -w Omit the DWARF symbol table.

Doc.: [https://pkg.go.dev/cmd/link](https://pkg.go.dev/cmd/link)

```bash
go build -ldflags="-s -w"
```
