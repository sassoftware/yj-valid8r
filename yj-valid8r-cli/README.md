# YJ-Valid8r: CLI

## Run via Direct Flags

```bash
go run main.go --schemas=path/to/schema.json --data=path/to/input.yaml
```

## Run via Config File

You can also run using a `config.yaml` file. Example:

```yaml
cliOutputFormat: "pretty" # Options: json | yaml | legacy | pretty (default)
strictValidation: true
checkTrailingWhitespace: true
schemas:
  - examples/schema.json
  - https://kubernetesjsonschema.dev/v1.10.3-standalone/service-v1.json
regexPatternRules:
  - name: Find Regex Pattern ${ }
    regex: '${(w+)(?::-[^}]*)?}'
searchPaths:
  - pathName: Get Target Ports
    pathKey: spec.ports[].targetPort
data: examples/data.yaml # YAML or JSON
```

Run the validator using the config file:

```bash
go run main.go --config=examples/config.yaml
```

## Override Config with Flags

You can override config values using command-line flags:

```bash
go run main.go --config=examples/config.yaml --schemas=path/to/schema.json --data=path/to/input.yaml
```

## Show Help / List All Flags

To display all available flags:

```bash
go run main.go --help
```