# intervals-api

## Getting started

### Initialize the project

```
go mod tidy
```

### Run the project

```
go run main.go
```

### Test the API

```
curl -X POST -H "Content-Type: application/json" -d '{"includes": [{"Start": 1, "End": 5}, {"Start": 9, "End": 13}], "excludes": [{"Start": 3, "End": 7}]}' http://localhost:8080/api/process
```

```
{"output":[{"Start":1,"End":2},{"Start":9,"End":13}]}
```