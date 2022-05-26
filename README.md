# WebShell

A simple WebShell service for basic testing and verification purposes.

```
go run main.go
# CGO_ENABLED=0 go build -o webshell main.go

curl http://localhost:8080/ping
curl -X POST http://localhost:8080/v1/exec -F cmd="ls -al /tmp"
curl -X POST http://localhost:8080/v1/exec -F cmd="whatever cmd you want to run"
```
