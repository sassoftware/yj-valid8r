# YJ-Valid8r: Web Playground

## Start the Web Server

```bash
go run main.go
```

By default, the web server starts at `http://localhost:7070`.

## Playground

You can access the interactive playground via your browser at:
`http://localhost:7070`

## API Usage

### Validate JSON via API

```bash
curl --request POST \
  --url http://localhost:7070/api/validate \
  --header 'Content-Type: application/json' \
  --data '{
    "schemas": [
      "/path/to/examples/schemas/schema.json"
    ],
    "data": "{\"id\":\"EMP-1001\",\"name\":\"Alice Johnson\",\"age\":35,\"role\":\"manager\",\"email\":\"alice@example.com\"}"
  }'
```

### Validate YAML via API

```bash
curl --request POST \
  --url http://localhost:7070/api/validate \
  --header 'Content-Type: application/json' \
  --data '{
    "schemas": [
      "/path/to/examples/schemas/schema.json"
    ],
    "data": "id: EMP-10A1\nname: \"\"\nage: 17\nrole: manager\nemail: \"aliceatexample.com\""
  }'
```