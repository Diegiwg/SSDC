# Starship Stop in Distance Calculator

This command-line interface (CLI) allows you to calculate the number of stops required for a starship to travel a given distance based on its Maximum Megalights per Hour (MGLT) value.You can test the calculations using either local or remote data sources.

## Build

To build the CLI, simply clone the repository and build the executable:

```bash
git clone https://github.com/Diegiwg/SSDC.git
cd SSDC
go build
```

## Usage

The CLI provides two main commands: `web`, `local` and `remote`.

### Docker Usage

You can use Docker to run the CLI:

```bash
docker run -p 8081:8081 diegiwg/ssdc:latest [command] [arguments]
```

Replace [command] and [arguments] with the appropriate CLI command and its arguments.

### Web Command

The `web` command allows you to test calculations using a website.

```bash
./SSDC[.exe] web [--local]
```

- `--local`: An optional flag to use local data.

### Local Command

The `local` command allows you to test calculations using locally stored data.

```bash
./SSDC[.exe] local <distance: int as MGLT> [--save: optional]
```

- `<distance>`: The distance to be traveled by the starship in Megalights.
- `--save`: An optional flag to save the data of API locally.

### Remote Command

The `remote` command enables testing with data from remote sources.

```bash
./SSDC[.exe] remote <distance: int as MGLT> [--save: optional]
```

- `<distance>`: The distance to be traveled by the starship in Megalights.
- `--save`: An optional flag to save the data of API locally.

## Examples

To calculate the number of stops needed for a starship to travel 100000 MGLT using local data:

```bash
./SSDC[.exe] local 100000
```

To calculate the same but with remote data and save the data locally:

```bash
./SSDC[.exe] remote 100000 --save
```

## Contributing

Contributions to this project are welcome. Feel free to submit bug reports, feature requests, or pull requests via the repository's issue tracker and pull request system.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
