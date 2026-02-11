# Contributing

If you're interested in contributing to this project, this is the best place to start. Before contributing to this project, please take a bit of time to read our [Code of Conduct](CODE-OF-CONDUCT.md). Also, note that this project is open-source and licensed under [Apache License 2.0](LICENSE).

## Project Structure

This project is a Go library providing a unified database abstraction layer.

- [db.go](db.go): Main interface definition.
- [query.go](query.go): Query builder.
- [model.go](model.go): Hook and helper interfaces.
- [driver/](driver/): Database driver implementations.

## Development

First prepare the backend environment by downloading all required dependencies:

```bash
go mod download
```

### Running Tests

You can run the full test suite (including all registered drivers) using:

```bash
go test ./... -v
```

## Adding New Drivers

New drivers (e.g., LevelDB, MongoDB) can be added by:
1. Creating a new directory under `driver/`.
2. Implementing the `dbase.Database` interface.
3. Registering the driver in an `init()` function via `dbase.Register`.
4. Adding the driver to `drivers/drivers.go` for convenience.

## Translations

If you would like to contribute to the documentation (like `README_CN.md`), feel free to submit a Pull Request.
