## Chapter 1 - Build a Commit Log Prototype
- A commit log is a data structure for an append-only sequence of records, ordered by time

### Testing

```bash
go run cmd/server/main.go

# produce
curl -X POST localhost:8080 -d \
    '{"record": {"value": "TGV0J3MgR28gIzEK"}}'

# consume
curl -X GET localhost:8080 -d '{"offset": 0}'
# {"Record":{"value":"TGV0J3MgR28gIzEK","offset":0}}
```

### Summary
- built a simple JSON/HTTP commit log service that accepts and responds with JSON
- store records in HTTP request body to in-memory log
